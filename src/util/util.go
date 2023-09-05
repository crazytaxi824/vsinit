package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// print color
const (
	COLOR_RED    = "\x1b[31m"
	COLOR_GREEN  = "\x1b[32m"
	COLOR_YELLOW = "\x1b[33m"

	COLOR_BOLD_RED    = "\x1b[1;31m"
	COLOR_BOLD_GREEN  = "\x1b[1;32m"
	COLOR_BOLD_YELLOW = "\x1b[1;33m"

	COLOR_RESET = "\x1b[0m"
)

// 询问是否在当前文件夹初始化项目
func Prompt(lang string) error {
	pwd, err := os.Getwd() // 获取当前路径
	if err != nil {
		return errors.New(projectCanceled(lang))
	}

	fmt.Printf("Init %s Project at %s%q%s? [Yes/no]: ",
		COLOR_BOLD_YELLOW+lang+COLOR_RESET, COLOR_YELLOW, pwd, COLOR_RESET)

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return errors.New(projectCanceled(lang))
	}

	if input != "yes\n" && input != "Yes\n" {
		return errors.New(projectCanceled(lang))
	}

	return nil
}

// 打印红色 cancel 信息
func projectCanceled(lang string) string {
	return fmt.Sprintf("%sInit %s Project Canceled!%s", COLOR_RED, lang, COLOR_RESET)
}
