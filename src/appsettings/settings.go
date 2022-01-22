package appsettings

import (
	"log"
	"os/user"
	"strings"
	"sync"
)

// ESLint 设置
const (
	// eslint global install path without filename. dir 用 '/' 结尾
	eslintGlobalPath = "~/.config/lints/"

	// javascript eslint filename
	JSESLintFileName = "eslintrc-js.json" // default: .eslintrc.json

	// typescript eslint filename
	TSESLintFileName = "eslintrc-ts.json" // default: .eslintrc.json

	// react eslint filename
	ReactESLintFileName = "eslintrc-react.json" // default: .eslintrc.json
)

var ESLintGlobalPath = defaultParsePath.parseESLintGlobalPath()

// 使用 sync.Once 只解析一次 filepath
type parsePath struct {
	value string
	s     sync.Once
}

// 实例化 parsePath
var defaultParsePath parsePath

func (p *parsePath) parseESLintGlobalPath() string {
	p.s.Do(func() {
		ps := strings.Split(eslintGlobalPath, "/")

		switch ps[0] {
		case "~":
			u, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}

			p.value = strings.Replace(eslintGlobalPath, "~", u.HomeDir, 1)

		case "":
			p.value = eslintGlobalPath

		default:
			log.Fatal("can not parse ESLintGlobalPath")
		}
	})

	return p.value
}
