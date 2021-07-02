package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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
		return "download vscode at https://code.visualstudio.com, or run:\nbrew install --cask visual-studio-code"

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

// 以下是检查 vscode extentsions
var extSuggestion = Suggestion{
	Problem:  "need to install following vscode extensions:",
	Solution: "you can install it in the vscode extentsion market, or run:\n",
}

// 检查 vscode 和 vscode 插件
func CheckVscodeAndExtensions(exts []string) ([]*Suggestion, error) {
	var suggs []*Suggestion

	sug := CheckCMDInstall("code")
	if sug != nil {
		// vscode 不存在
		suggs = append(suggs, sug, &Suggestion{
			Problem: extSuggestion.Problem,
			Solution: extSuggestion.Solution + "code --install-extension " +
				strings.Join(exts, "; \\\ncode --install-extension "),
		})
	} else {
		// vscode 存在
		cmd := exec.Command("code", "--list-extensions")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		var result []string

		for _, ext := range exts {
			if !bytes.Contains(out, []byte(ext)) {
				result = append(result, "code --install-extension "+ext)
			}
		}

		if len(result) > 0 {
			suggs = append(suggs, &Suggestion{
				Problem:  extSuggestion.Problem,
				Solution: extSuggestion.Solution + strings.Join(result, "; \\\n"),
			})
		}
	}

	if len(suggs) > 0 {
		return suggs, nil
	}

	return nil, nil
}
