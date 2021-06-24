// vsc envcheck go -lint

package golang

import (
	"bytes"
	"errors"
	"local/src/util"
	"os"
	"os/exec"
)

func CheckGO(lintFlag bool) ([]*util.Suggestion, error) {
	return checkGOENV(lintFlag)
}

// $GOBIN 是否存在
func checkGOENV(lintFlag bool) ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 SHELL 环境设置
	sug := checkGOPATH()
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 go 是否安装 // 检查 code 安装
	sug = util.CheckCMDInstall("go")
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 vscode and extensions,
	sug = util.CheckCMDInstall("code")
	if sug != nil {
		suggs = append(suggs, sug, &util.Suggestion{ // 安装 vscode 插件 GO
			Problem: "need to install vscode extension 'golang.go'",
			Solution: "you can install it in the vscode extentsion market, or run\n" +
				"code --install-extension golang.go",
		})
	} else {
		su, er := checkVscodeExtensions()
		if er != nil {
			return nil, er
		}
		if su != nil {
			suggs = append(suggs, su)
		}
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
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}

// 安装 vscode 插件 GO
func checkVscodeExtensions() (*util.Suggestion, error) {
	cmd := exec.Command("code", "--list-extensions")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if !bytes.Contains(out, []byte("golang.go")) {
		return &util.Suggestion{
			Problem: "need to install vscode extension 'golang.go'",
			Solution: "you can install it in the vscode extentsion market, or run\n" +
				"code --install-extension golang.go",
		}, nil
	}
	return nil, nil
}

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

// TODO
func checkGolangciLint() (*util.Suggestion, error) {
	cfg, suggestion, err := util.ReadVscFile()
	if err != nil {
		return nil, err
	}
	if suggestion != nil {
		return suggestion, nil
	}

	if cfg.Golangci == "" {
		return &util.Suggestion{
			Problem:  "haven't setup golangci-lint yet, please set it:",
			Solution: util.GolintciCmd,
		}, nil
	}

	gof, err := os.Open(cfg.Golangci)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		return &util.Suggestion{
			Problem:  "golangci-lint config file is missing, please re-install it:",
			Solution: util.GolintciCmd,
		}, nil
	}
	defer gof.Close()

	// 能够打开说明已经设置成功
	return nil, nil
}
