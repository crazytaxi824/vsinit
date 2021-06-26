package util

import (
	"bytes"
	"errors"
)

func isComment(src []byte) bool {
	tmp := bytes.TrimSpace(src)
	l := len(tmp)
	if l == 0 { // 空行, 算作 comment
		return true
	} else if l == 1 { // {} [] , 等情况
		return false
	}

	if string(tmp[:2]) == "//" {
		return true
	}

	return false
}

func analyseJSONCstatement(src []byte) (lastCharIndex int, hasComments bool, err error) {
	l := len(src)

	// NOTE lastIndex = -1, 说明该行是空行。
	var lastIndex = -1  // 最后一位有效 char 的 index。
	var quote bool      // 在引号内还是引号外
	var transfer bool   // 是否在转义状态
	var hasComment bool // 是否有 comment
	// var multiLineCommentMark bool // FIXME 多行注释 /* */

	// 逐字判断
	for i := 0; i < l; i++ {
		if transfer { // 如果转义了，忽略后面一个char
			transfer = false
			lastIndex = i
			continue
		}

		if src[i] == ' ' {
			// 如果是空则不移动 lastIndex
			continue
		} else if src[i] == '"' {
			// 标记是 opening quote 还是 closing quote.
			if quote {
				quote = false
			} else {
				quote = true
			}
			lastIndex = i // 移动 lastIndex
		} else if src[i] == '\\' {
			// 如果是转义符，则标记转义
			if !quote {
				// 如果转义符在引号外面，ERROR
				return 0, false, errors.New("error: '\\' out side quote")
			}
			transfer = true
		} else if src[i] == '/' {
			if quote {
				// 如果 ‘/’ 在引号内，不需要特殊处理
				lastIndex = i
				continue
			}

			// 如果 '/' 在引号外面，判断后一位是否也是 '/'，说明后面是 comments.
			if i+1 < l && src[i+1] == '/' {
				hasComment = true // 标记 comments
				break             // 结束循环
			}
			// FIXME /* */ 多行注释问题

			// 如果 ‘/’ 在引号外面而且后面不是 ‘/’ ，ERROR
			return 0, false, errors.New("error: '/' out side quote")
		} else {
			// 其他正常情况下直接向后处理。
			lastIndex = i
		}
	}

	// 如果 line 结束，单引号没有关闭则，Error
	if quote {
		return 0, false, errors.New("error: statement is Unquoted")
	}

	return lastIndex, hasComment, nil
}

// NOTE JSONC must be formatted, otherwise cannot be read.
func JSONCToJSON(jsonc []byte) ([]byte, error) {
	lines := bytes.Split(jsonc, []byte("\n"))

	var result [][]byte
	for _, line := range lines {
		if isComment(line) {
			continue
		}

		lastIndex, _, er := analyseJSONCstatement(line)
		if er != nil {
			return nil, er
		}
		result = append(result, line[:lastIndex+1])
	}

	return bytes.Join(result, []byte("\n")), nil
}

// TODO 插入数据到 JSONC 中
