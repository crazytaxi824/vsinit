package js

import (
	"bytes"
	"local/src/util"
	"os/exec"
)

var prettierSuggestion = util.Suggestion{
	Problem: "need to install vscode extension 'esbenp.prettier-vscode'",
	Solution: "you can install it in the vscode extentsion market, or run:\n" +
		"code --install-extension esbenp.prettier-vscode",
}

func checkJS() ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 go 是否安装 // 检查 code 安装
	sug := util.CheckCMDInstall("node", "npm")
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 vscode and extensions,
	sug = util.CheckCMDInstall("code")
	if sug != nil {
		suggs = append(suggs, sug, &prettierSuggestion)
	} else {
		// 已经安装了 vscode
		su, er := checkVscodeExtensions()
		if er != nil {
			return nil, er
		}
		if su != nil {
			suggs = append(suggs, su)
		}
	}

	// 检查返回是否为空
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}

// 检查 vscode 插件 prettier
func checkVscodeExtensions() (*util.Suggestion, error) {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if !bytes.Contains(out, []byte("esbenp.prettier-vscode")) {
		return &prettierSuggestion, nil
	}
	return nil, nil
}

// 检查是否安装了 jest
func checkJest() *util.Suggestion {
	return util.CheckCMDInstall("jest")
}

// 检查是否安装了 eslint
func checkESLint() *util.Suggestion {
	return util.CheckCMDInstall("eslint")
}
