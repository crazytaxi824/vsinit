// vsc envcheck go -lint

package golang

import (
	"fmt"
	"local/src/util"
	"os"
)

// $GOBIN 是否存在
func checkGOENV() error {
	// 检查 go 是否安装 // 检查 code 安装
	err := util.CheckCMDInstall("go", "code")
	if err != nil {
		return err
	}
	// 检查 SHELL 环境设置
	err = checkGOPATH()
	if err != nil {
		return err
	}

	// 检查 vscode extension,

	// plugins:gopkgs,go-outline,gotests,gomodifytags,impl,dlv,golangci-lint,gopls
	// 检查 vscode extension 工具链.
	// go get xxxx 安装.

	// 检查 settings，launch 设置.
	return nil
}

// TODO 检查 lint config file 位置。
func checkGolangciLint() {

}

func checkGOPATH() error {
	if os.Getenv("GOPATH") == "" {
		// return errors.New("please setup $GOPATH in your environment, ~/.bash_profile OR ./zshrc")
		return fmt.Errorf(`please setup $GOPATH in your environment, ~/.bash_profile OR ./zshrc.
'''
# golang setting
# export GOROOT=/usr/local/go  # 'echo $PATH' should have $GOROOT/bin
export GOPATH=/Users/ray/gopath
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN
export GO111MODULE=on
'''`)
	}
	return nil
}
