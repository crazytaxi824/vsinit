package util

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// 询问是否在当前文件夹初始化项目
func Prompt(lang string) error {
	pwd, err := os.Getwd() // 获取当前路径
	if err != nil {
		return fmt.Errorf("%s: %w", Cancel(), err)
	}

	fmt.Printf("Init %s Project at %s%q%s? [y(es)/no]: ",
		COLOR_BOLD_YELLOW+lang+COLOR_RESET, COLOR_YELLOW, pwd, COLOR_RESET)

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return fmt.Errorf("%s: %w", Cancel(), err)
	}
	input = strings.ToLower(strings.TrimSpace(input))

	if input != "y" && input != "yes" {
		return errors.New(Cancel())
	}

	return nil
}

// 打印红色 cancel 信息
func Cancel() string {
	return fmt.Sprintf("%sCancel%s", COLOR_BOLD_RED, COLOR_RESET)
}
