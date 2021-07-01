package golang

import (
	"local/src/util"
	"os"
	"strings"
)

// 需要安装的插件
var extensions = []string{"golang.go", "humao.rest-client"}

// go 插件所需 tools
var goTools = []string{"gopkgs", "go-outline", "gotests",
	"gomodifytags", "impl", "dlv", "golangci-lint", "gopls"}

func CheckGO() ([]*util.Suggestion, error) {
	return checkGOENV()
}

// 检查所有 GO 运行环境
func checkGOENV() ([]*util.Suggestion, error) {
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

	// goTools: gopkgs,go-outline,gotests,gomodifytags,impl,dlv,golangci-lint,gopls
	// 检查 vscode extension 工具链.
	// go get xxxx 安装.
	sug = checkGoTools(goTools...)
	if sug != nil {
		suggs = append(suggs, sug)
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
				"export GOPATH=/Users/ray/gopath\n" + // FIXME
				"export GOBIN=$GOPATH/bin\n" +
				"export PATH=$PATH:$GOBIN\n" +
				"export GO111MODULE=on",
		}
	}
	return nil
}

// 检查 gotools 安装情况
// goTools: gopkgs,go-outline,gotests,gomodifytags,impl,dlv,golangci-lint,gopls
func checkGoTools(tools ...string) *util.Suggestion {
	var solutions []string
	for _, tool := range tools {
		if !util.CheckCommandExist(tool) {
			solutions = append(solutions, goToolsSuggestion(tool))
		}
	}

	if len(solutions) > 0 {
		return &util.Suggestion{
			Problem:  "need to install following goTools:",
			Solution: strings.Join(solutions, "; \\\n"),
		}
	}

	return nil
}

// 检查 vscode 中 go 插件所需要的工具.
func goToolsSuggestion(tool string) string {
	switch tool {
	case "gopkgs":
		return "go get github.com/uudashr/gopkgs/v2/cmd/gopkgs"
	case "go-outline":
		return "go get github.com/ramya-rao-a/go-outline"
	case "gotests":
		return "go get github.com/cweill/gotests/gotests"
	case "impl":
		return "go get github.com/josharian/impl"
	case "dlv":
		return "go get github.com/go-delve/delve/cmd/dlv"
	case "gopls":
		return "go get golang.org/x/tools/gopls"
	case "golangci-lint":
		return "go get github.com/golangci/golangci-lint/cmd/golangci-lint"
	case "gomodifytags":
		return "go get github.com/fatih/gomodifytags"
	case "debug-cmd":
		return "this is a debug test solution."
	}
	return util.InternalErrMsg
}
