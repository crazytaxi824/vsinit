// vsc envcheck go -lint

package golang

import (
	"errors"
	"local/src/util"
	"os"
)

// FIXME
const GolintciCmd = "vs init go -cilint <path>"

// 需要安装的插件
var extensions = []string{"golang.go", "humao.rest-client"}

func CheckGO(lintFlag bool) ([]*util.Suggestion, error) {
	return checkGOENV(lintFlag)
}

// 检查所有 GO 运行环境
func checkGOENV(lintFlag bool) ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 SHELL 环境设置
	sug := checkGOPATH()
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 go 是否安装
	sug = util.CheckCMDInstall("go")
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

	// plugins:gopkgs,go-outline,gotests,gomodifytags,impl,dlv,golangci-lint,gopls
	// 检查 vscode extension 工具链.
	// go get xxxx 安装.
	sug = util.CheckCMDInstall(util.GoTools...)
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 golangci-lint
	if lintFlag {
		su, er := checkGolangciLint()
		if er != nil {
			return nil, er
		}
		if su != nil {
			suggs = append(suggs, su)
		}
	}

	// 检查返回是否为空
	if len(suggs) > 0 {
		return suggs, nil
	}

	return nil, nil
}

// 检查 $GOPATH 设置
func checkGOPATH() *util.Suggestion {
	if os.Getenv("GOPATH") == "" {
		return &util.Suggestion{
			Problem: "need to setup $GOPATH in ~/.bash_profile OR ./zshrc",
			Solution: "# golang setting\n" +
				"export GOPATH=/Users/ray/gopath\n" +
				"export GOBIN=$GOPATH/bin\n" +
				"export PATH=$PATH:$GOBIN\n" +
				"export GO111MODULE=on",
		}
	}
	return nil
}

// 检查 golang-ci lint 设置
func checkGolangciLint() (*util.Suggestion, error) {
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return nil, err
	}

	// 读取 vsc setting
	var vscCfgJSON util.VscConfigJSON
	err = vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		return &util.Suggestion{
			Problem:  "haven't setup golangci-lint yet, please run:",
			Solution: GolintciCmd,
		}, nil
	}

	// 查找 golangci 设置
	if vscCfgJSON.Golangci == "" {
		return &util.Suggestion{
			Problem:  "haven't setup golangci-lint yet, please run:",
			Solution: GolintciCmd,
		}, nil
	}

	// 寻找 golangci 配置文件
	gof, err := os.Open(vscCfgJSON.Golangci)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		return &util.Suggestion{
			Problem:  "golangci-lint config file is missing, please run:",
			Solution: GolintciCmd,
		}, nil
	}
	defer gof.Close()

	// 能够打开说明已经设置成功
	return nil, nil
}
