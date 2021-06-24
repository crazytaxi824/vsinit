package util

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

var GoTools = []string{"gopkgs", "go-outline", "gotests",
	"gomodifytags", "impl", "dlv", "golangci-lint", "gopls"}

type ErrorMsg struct {
	Problem  string
	Solution []string
}

func (e ErrorMsg) Error() string {
	return Warn(">>> "+e.Problem) + "\n" + strings.Join(e.Solution, "\n") + "\n\n"
}

type Erros []error

func (es Erros) Error() string {
	var builder strings.Builder
	for _, err := range es {
		builder.WriteString(err.Error())
	}
	return builder.String()
}

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
		solution := installMsg(v)
		solutions = append(solutions, solution...)
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

func installMsg(cmdName string) []string {
	switch cmdName {
	case "code":
		return []string{"you can download it at https://code.visualstudio.com"}

	case "go":
		return []string{"you can download it at https://golang.org/"}

	case "node", "npm":
		return []string{"you can download it at https://nodejs.org/"}

	case "tsc":
		return []string{"you can run 'npm i -g typescript' at terminal"}

	case "jest":
		return []string{"you can run 'npm i -g jest' at terminal"}

	case "eslint":
		return []string{"you can run 'npm i -g eslint' at terminal"}

	case "debug-cmd", "gopkgs", "go-outline", "gotests", "gomodifytags", "impl", "dlv", "golangci-lint", "gopls":
		// DEBUG
		return checkGoTools(cmdName)
	}

	return []string{"please contact author"}
}

func checkGoTools(tool string) []string {
	var solutions []string

	switch tool {
	case "gopkgs":
		solutions = append(solutions, "go get github.com/uudashr/gopkgs/v2/cmd/gopkgs")
	case "go-outline":
		solutions = append(solutions, "go get github.com/ramya-rao-a/go-outline")
	case "gotests":
		solutions = append(solutions, "go get github.com/cweill/gotests/gotests")
	case "impl":
		solutions = append(solutions, "go get github.com/josharian/impl")
	case "dlv":
		solutions = append(solutions, "go get github.com/go-delve/delve/cmd/dlv")
	case "gopls":
		solutions = append(solutions, "go get golang.org/x/tools/gopls")
	case "golangci-lint":
		solutions = append(solutions, "go get github.com/golangci/golangci-lint/cmd/golangci-lint")
	case "gomodifytags":
		solutions = append(solutions, "go get github.com/fatih/gomodifytags")
		// case "debug-cmd":
		// 	solutions = append(solutions, "this is a debug test solution.")
	}

	return solutions
}
