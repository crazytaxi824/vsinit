package ts

import (
	"local/src/util"
)

// 需要安装的插件
var extensions = []string{"esbenp.prettier-vscode",
	"VisualStudioExptTeam.vscodeintellicode",
	"christian-kohler.path-intellisense",
	"dbaeumer.vscode-eslint",
}

func CheckTS() ([]*util.Suggestion, error) {
	return checkTS()
}

func checkTS() ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 node 是否安装
	sug := util.CheckCMDInstall("node")
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 typescript 是否安装
	sug = util.CheckCMDInstall("tsc")
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
	if len(suggs) > 0 {
		return suggs, nil
	}

	return nil, nil
}
