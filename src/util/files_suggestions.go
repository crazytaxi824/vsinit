package util

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// 需要写入项目的文件
type FileContent struct {
	Dir        string // 文件夹地址, 用于 mkdir
	FileName   string // 文件名, 不带路径
	Content    []byte // 文件内容
	Overwrite  bool   // 是否需要覆盖文件内容
	Suggestion string // 如果文件冲突, 提供建议
}

func WriteAllFiles(fileContents []FileContent) error {
	// 检查 dir 是否以 "/" 结尾
	err := checkDir(fileContents)
	if err != nil {
		return err
	}

	// create dir and write files
	err = createDirAndWriteFiles(fileContents)
	if err != nil {
		return err
	}

	return nil
}

// 写文件之前先检查 dir 是否以 "/" 结尾
func checkDir(fileContents []FileContent) error {
	for _, fc := range fileContents {
		if fc.Dir != "" && fc.Dir[len(fc.Dir)-1] != '/' {
			errmsg := fmt.Sprintf("%s dir is not ended with '/'", fc.Dir)
			log.Println(errmsg)
			return errors.New(errmsg)
		}
	}

	return nil
}

// 写文件
func createDirAndWriteFiles(fileContents []FileContent) error {
	fmt.Printf("%sWriting files ...%s\n", COLOR_GREEN, COLOR_RESET)

	// create dir, if Dir == "" skip
	for _, fc := range fileContents {
		if fc.Dir != "" {
			err := os.MkdirAll(fc.Dir, 0750)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		// write files
		var (
			f   *os.File
			err error
		)

		// check overwrite option
		if fc.Overwrite {
			f, err = os.OpenFile(fc.Dir+fc.FileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		} else {
			f, err = os.OpenFile(fc.Dir+fc.FileName, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
		}
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Println(err)
			return err
		} else if errors.Is(err, os.ErrExist) {
			fmt.Printf(" - %s ... skip\n", fc.Dir+fc.FileName)

			// [VVI] 如果文件存在检查 suggestion.
			if fc.Suggestion != "" {
				suggestions = append(suggestions, Suggestion{
					FileName: fc.Dir + fc.FileName,
					Msg:      fc.Suggestion,
				})
			}
			continue
		}

		_, err = f.Write(fc.Content)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Printf("%s - %s ... done%s\n", COLOR_GREEN, fc.Dir+fc.FileName, COLOR_RESET)
	}

	return nil
}

// files 中的 suggestions 放到这里统一打印在最后.
var suggestions []Suggestion

type Suggestion struct {
	FileName string // 建议修改的文件
	Msg      string // 建议内容
}

func PrintSuggestions() {
	fmt.Printf("\n%syou might need to change following settings manually:%s\n\n", COLOR_BOLD_YELLOW, COLOR_RESET)
	for _, sug := range suggestions {
		fmt.Printf(">>> %s%q:%s <<<\n", COLOR_GREEN, sug.FileName, COLOR_RESET)
		fmt.Println(sug.Msg)
	}
}
