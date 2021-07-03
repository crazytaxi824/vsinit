package js

import (
	"local/src/util"
)

// 需要安装的 vscode 插件
var extensions = []string{"esbenp.prettier-vscode",
	"VisualStudioExptTeam.vscodeintellicode",
	"christian-kohler.path-intellisense",
	"dbaeumer.vscode-eslint",
}

func CheckJS() ([]*util.Suggestion, error) {
	return checkJS()
}

func checkJS() ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 node, typescript 是否安装
	sug := util.CheckCMDInstall("node")
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 vscode and extensions,
	sus, err := util.CheckVscodeAndExtensions(extensions)
	if err != nil {
		return nil, err
	}
	if len(sus) > 0 {
		suggs = append(suggs, sus...)
	}

	// 检查返回是否为空
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}
