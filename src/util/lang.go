package util

import (
	"os/exec"
	"runtime"
	"strings"
)

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
	var es Erros
	for _, lang := range langs {
		if err := checkCommandExistence(lang); err != nil {
			es = append(es, err)
		}
	}

	if len(es) == 0 {
		return nil
	}
	return es
}

// 'which <cmd>'
func checkCommandExistence(cmdName string) error {
	cmd := exec.Command(whichCmd(), cmdName)
	err := cmd.Run()
	if err != nil {
		return installMsg(cmdName)
	}
	return nil
}

// linux & mac(darwin) using which, windows using where
func whichCmd() string {
	if runtime.GOOS == "windows" {
		return "where"
	}
	return "which"
}

func installMsg(cmdName string) ErrorMsg {
	switch cmdName {
	case "code":
		return ErrorMsg{
			Problem:  "please install 'vscode' first",
			Solution: []string{"you can download it at https://code.visualstudio.com"},
		}

	case "go":
		return ErrorMsg{
			Problem:  "please install 'go' first",
			Solution: []string{"you can download it at https://golang.org/"},
		}

	case "node", "npm":
		return ErrorMsg{
			Problem:  "please install 'nodejs' first",
			Solution: []string{"you can download it at https://nodejs.org/"},
		}

	case "tsc":
		return ErrorMsg{
			Problem:  "please install 'typescript' first",
			Solution: []string{"you can run 'npm i -g typescript' at terminal"},
		}

	case "jest":
		return ErrorMsg{
			Problem:  "please install 'jest' first",
			Solution: []string{"you can run 'npm i -g jest' at terminal"},
		}

	case "eslint":
		return ErrorMsg{
			Problem:  "please install 'eslint' first",
			Solution: []string{"you can run 'npm i -g eslint' at terminal"},
		}

	case "gopkgs", "go-outline", "gotests", "gomodifytags", "impl", "dlv", "golangci-lint", "gopls":
		return ErrorMsg{
			Problem:  "please install following tools:",
			Solution: checkGoTools(cmdName),
		}
	}

	return ErrorMsg{
		Problem:  "command is not in the list",
		Solution: []string{"please contact author"},
	}
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
	}

	return solutions
}
