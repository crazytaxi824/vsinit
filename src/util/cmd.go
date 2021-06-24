package util

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var GoTools = []string{"gopkgs", "go-outline", "gotests",
	"gomodifytags", "impl", "dlv", "golangci-lint", "gopls"}

// 检查是否安装了语言
func CheckCMDInstall(langs ...string) *Suggestion {
	var result []string
	for _, lang := range langs {
		if !checkCommandExist(lang) {
			result = append(result, lang)
		}
	}

	if len(result) == 0 {
		return nil
	}

	var solutions []string
	for _, v := range result {
		solut := installMsg(v)
		solutions = append(solutions, solut)
	}

	return &Suggestion{
		Problem:  fmt.Sprintf("need to intall '%s':", strings.Join(result, ", ")),
		Solution: strings.Join(solutions, "; \\\n"),
	}
}

// 'which <cmd>'
func checkCommandExist(cmdName string) bool {
	cmd := exec.Command(whichCmd(), cmdName)
	err := cmd.Run()
	return err == nil
}

// linux & mac(darwin) using which, windows using where
func whichCmd() string {
	if runtime.GOOS == "windows" {
		return "where"
	}
	return "which"
}

func installMsg(cmdName string) string {
	switch cmdName {
	case "code":
		return "download it at https://code.visualstudio.com, or run:\nbrew install vscode"

	case "go":
		return "download it at https://golang.org/, or run:\nbrew install go"

	case "node", "npm":
		return "download it at https://nodejs.org/, or run:\nbrew install node"

	case "tsc":
		return "npm i -g typescript"

	case "jest":
		return "npm i -g jest"

	case "eslint":
		return "npm i -g eslint"

	case "debug-cmd", "gopkgs", "go-outline", "gotests", "gomodifytags", "impl", "dlv", "golangci-lint", "gopls":
		// DEBUG
		return checkGoTools(cmdName)
	}

	return InternalErrMsg
}

func checkGoTools(tool string) string {
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
	return InternalErrMsg
}
