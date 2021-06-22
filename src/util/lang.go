package util

import (
	"errors"
	"os/exec"
	"runtime"
)

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
	case "golangci-lint":
		return errors.New("please install 'golangci-lint', https://golangci-lint.run/")
	case "node", "npm":
		return errors.New("please install 'node' first, https://nodejs.org/")
	case "tsc":
		return errors.New("please install 'typescript' first, 'npm i -g typescript'")
	case "jest":
		return errors.New("please install 'jest' first, 'npm i -g jest'")
	case "eslint":
		return errors.New("please install 'eslint' first, 'npm i -g eslint'")
	case "python", "python3":
		return errors.New(`please install 'python' first, https://www.python.org, and then
		install python extension 'code --install-extension ms-python.python'`)
	}
	return nil
}
