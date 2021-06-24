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
func CheckCMDInstall(langs ...string) error {
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

	return ErrorMsg{
		Problem:  fmt.Sprintf("please intall '%s' first:", strings.Join(result, ", ")),
		Solution: solutions,
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
		return "you can download it at https://code.visualstudio.com"

	case "go":
		return "you can download it at https://golang.org/"

	case "node", "npm":
		return "you can download it at https://nodejs.org/"

	case "tsc":
		return "you can run 'npm i -g typescript' at terminal"

	case "jest":
		return "you can run 'npm i -g jest' at terminal"

	case "eslint":
		return "you can run 'npm i -g eslint' at terminal"

	case "debug-cmd", "gopkgs", "go-outline", "gotests", "gomodifytags", "impl", "dlv", "golangci-lint", "gopls":
		// DEBUG
		return checkGoTools(cmdName)
	}

	return ErrInternalMsg
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
	return ErrInternalMsg
}
