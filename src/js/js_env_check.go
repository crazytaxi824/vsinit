package js

import (
	"local/src/util"
	"os"
)

var extensions = []string{"esbenp.prettier-vscode",
	"VisualStudioExptTeam.vscodeintellicode",
	"christian-kohler.path-intellisense",
	"dbaeumer.vscode-eslint",
}

func CheckJS(tsjs util.TSJSFlags) ([]*util.Suggestion, error) {
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjs.FlagSet.Parse(os.Args[3:])

	return checkJS(*tsjs.Jest, *tsjs.ESLint)
}

func checkJS(jest, eslint bool) ([]*util.Suggestion, error) {
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

	// TODO jest, eslint
	if jest {
		sug = checkJest()
		if sug != nil {
			suggs = append(suggs, sug)
		}
	}

	if eslint {
		sug = checkESLint()
		if sug != nil {
			suggs = append(suggs, sug)
		}
	}

	// 检查返回是否为空
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}

// 检查是否安装了 jest
func checkJest() *util.Suggestion {
	return util.CheckCMDInstall("jest")
}

// 检查是否安装了 eslint
func checkESLint() *util.Suggestion {
	return util.CheckCMDInstall("eslint")
}
