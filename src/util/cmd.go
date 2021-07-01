package util

import (
	"fmt"
	"os/exec"
)

// 检查是否安装了命令行工具
func CheckCMDInstall(lang string) *Suggestion {
	if !CheckCommandExist(lang) {
		return &Suggestion{
			Problem:  fmt.Sprintf("need to intall '%s':", lang),
			Solution: solutionMsg(lang),
		}
	}
	return nil
}

// 'which <cmd>'
func CheckCommandExist(cmdName string) bool {
	cmd := exec.Command("which", cmdName)
	err := cmd.Run()
	return err == nil
}

// 如果缺失以下 command line 工具，则提示以下内容.
func solutionMsg(cmdName string) string {
	switch cmdName {
	case "code":
		return "download vscode at https://code.visualstudio.com, or run:\nbrew install vscode"

	case "go":
		return "download go at https://golang.org/, or run:\nbrew install go"

	case "node", "npm":
		return "download node at https://nodejs.org/, or run:\nbrew install node"

	case "python", "python3":
		return "download python at https://www.python.org/, or run:\nbrew install python3"

	case "tsc":
		return "npm i -g typescript"

	case "jest":
		return "npm i -g jest"

	case "eslint":
		return "npm i -g eslint"
	}

	return InternalErrMsg
}
