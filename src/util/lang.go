package util

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

type ErrorMsg struct {
	Problem  string
	Solution []string
}

func (e ErrorMsg) Error() string {
	return Warn(">>>>>> "+e.Problem) + "\n" + strings.Join(e.Solution, "\n")
}

// 检查是否安装了语言
func CheckCMDInstall(langs ...string) error {
	for _, lang := range langs {
		if err := checkCommandExistence(lang); err != nil {
			return err
		}
	}
	return nil
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

func installMsg(cmdName string) error {
	switch cmdName {
	case "code":
		return errors.New(`please install 'VScode' first, 
		this Project init tool is base on VScode env.
		download it at https://code.visualstudio.com`)

	case "go":
		return errors.New(`please install 'go' first, https://golang.org/, and then
		install go extension 'code --install-extension golang.go'`)

	case "python", "python3":
		return errors.New(`please install 'python' first, https://www.python.org, and then
		install python extension 'code --install-extension ms-python.python'`)

	case "node", "npm":
		return errors.New("please install 'nodejs' first, https://nodejs.org/")

	case "tsc":
		return errors.New("please install 'typescript' first, 'npm i -g typescript'")

	case "jest":
		return errors.New("please install 'jest' first, 'npm i -g jest'")
	}

	return errors.New("command is not in the list, please contact author")
}
