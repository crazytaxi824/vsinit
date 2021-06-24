// vsc envcheck go -lint

package golang

import (
	"bytes"
	"local/src/util"
	"os"
	"os/exec"
)

// NOTE 需要和 util 中的 case match.
var goTools = []string{"gopkgs", "go-outline", "gotests", "gomodifytags", "impl", "dlv", "golangci-lint", "gopls"}

// $GOBIN 是否存在
func checkGOENV() error {
	var ers util.Erros

	// 检查 SHELL 环境设置
	err := checkGOPATH()
	if err != nil {
		ers = append(ers, err)
	}

	// 检查 go 是否安装 // 检查 code 安装
	err = util.CheckCMDInstall("go")
	if err != nil {
		ers = append(ers, err)
	}

	// 检查 vscode and extensions,
	err = util.CheckCMDInstall("code")
	if err != nil {
		ers = append(ers, err, util.ErrorMsg{
			Problem: "please install vscode extension 'golang.go'",
			Solution: []string{"you can install it in the vscode extentsion market,",
				"or run 'code --install-extension golang.go'"},
		})
	} else {
		er := checkVscodeExtensions()
		if er != nil {
			ers = append(ers, er)
		}
	}

	// plugins:gopkgs,go-outline,gotests,gomodifytags,impl,dlv,golangci-lint,gopls
	// 检查 vscode extension 工具链.
	// go get xxxx 安装.
	err = util.CheckCMDInstall(goTools...)
	if err != nil {
		ers = append(ers, err)
	}

	// TODO 检查 settings，launch 设置.

	// 检查返回是否为空
	if len(ers) == 0 {
		return nil
	}

	return ers
}

func checkVscodeExtensions() error {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	if !bytes.Contains(out, []byte("golang.go")) {
		return util.ErrorMsg{
			Problem: "please install vscode extension 'golang.go'",
			Solution: []string{"you can install it in the vscode extentsion market,",
				"or run 'code --install-extension golang.go'"},
		}
	}
	return nil
}

func checkGOPATH() error {
	if os.Getenv("GOPATH") == "" {
		return util.ErrorMsg{
			Problem: "please setup $GOPATH in your environment, ~/.bash_profile OR ./zshrc",
			Solution: []string{"```",
				"# golang setting",
				"export GOPATH=/Users/ray/gopath",
				"export GOBIN=$GOPATH/bin",
				"export PATH=$PATH:$GOBIN",
				"export GO111MODULE=on",
				"'''",
			},
		}
	}
	return nil
}

// // TODO 检查 lint config file 位置。
// func checkGolangciLint() {
// }
