package python

import (
	"local/src/util"
)

// 需要安装的 vscode 插件
var extensions = []string{"ms-python.python",
	"ms-python.vscode-pylance",
	"VisualStudioExptTeam.vscodeintellicode",
}

func CheckPython() ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 python3 是否安装
	sug := util.CheckCMDInstall("python3")
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 vscode 和 vscode extensions
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
