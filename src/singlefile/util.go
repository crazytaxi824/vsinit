package singlefile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"local/src/util"
)

// 获取命令行的前两个参数 - eg: "vs file xxx"
func getCmd() string {
	return fmt.Sprintf("%q", strings.Join(os.Args[:3], " "))
}

// 询问是否在当前文件夹初始化项目
func askBeforeProceed(lang string) error {
	pwd, err := os.Getwd() // 获取当前路径
	if err != nil {
		return errors.New(writeFileCanceled(lang))
	}

	fmt.Printf("Write file %s%q%s at %s%q%s? [Yes/no]: ",
		util.COLOR_BOLD_YELLOW, lang, util.COLOR_RESET, util.COLOR_YELLOW, pwd+"/", util.COLOR_RESET)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return errors.New(writeFileCanceled(lang))
	}

	if input != "yes\n" && input != "Yes\n" {
		return errors.New(writeFileCanceled(lang))
	}

	return nil
}

// 打印红色 cancel 信息
func writeFileCanceled(lang string) string {
	return fmt.Sprintf("%sWrite file %q Canceled!%s", util.COLOR_RED, lang, util.COLOR_RESET)
}
