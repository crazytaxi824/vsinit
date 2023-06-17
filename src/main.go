package main

import (
	"errors"
	"fmt"
	"os"

	"local/src/golang"
	"local/src/javascript"
	"local/src/react"
	"local/src/singlefile"
	"local/src/typescript"
	"local/src/util"
)

func main() {
	// 只输入了 'vs' 命令的情况下
	if len(os.Args) < 2 {
		fmt.Printf(mainHelpMsg, util.COLOR_GREEN, util.COLOR_RESET, util.COLOR_GREEN, util.COLOR_RESET)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "go":
		if err := golang.InitGoProj(); err != nil {
			os.Exit(2)
		}

	case "ts", "typescript":
		if err := typescript.InitTSProj(); err != nil {
			os.Exit(2)
		}

	case "js", "javascript":
		if err := javascript.InitJSProj(); err != nil {
			os.Exit(2)
		}

	case "react":
		if err := react.InitReactProj(); err != nil {
			os.Exit(2)
		}

	case "file":
		if err := writeFile(); err != nil {
			os.Exit(2)
		}

	default:
		fmt.Printf(mainHelpMsg, util.COLOR_GREEN, util.COLOR_RESET, util.COLOR_GREEN, util.COLOR_RESET)
		os.Exit(2)
	}

	fmt.Printf("%sAll Done! Happy Coding!%s\n", util.COLOR_BOLD_GREEN, util.COLOR_RESET)
}

const mainHelpMsg = `Init a project at current directory
Usage:%s
    vs [go | ts | js | react | file]%s

flags info:%s
    vs ts -h
    vs js -h
    vs react -h
    vs file -h%s

`

// write a specific file.
func writeFile() error {
	// 只输入了 'vs file' 命令的情况下
	if len(os.Args) < 3 {
		fmt.Printf(fileHelpMsg, util.COLOR_GREEN, util.COLOR_RESET, util.COLOR_GREEN, util.COLOR_RESET)
		return errors.New("filename not specified")
	}

	switch os.Args[2] {
	case "editorconfig":
		return singlefile.WriteEditorConfigFile()

	case "golangci":
		return singlefile.WriteGolangciFile()
	}

	fmt.Printf(fileHelpMsg, util.COLOR_GREEN, util.COLOR_RESET, util.COLOR_GREEN, util.COLOR_RESET)
	return errors.New("specified filename does not exist")
}

const fileHelpMsg = `write a single file at current directory
Usage:%s
    vs file [editorconfig | golangci]%s

flags info:%s
    vs file editorconfig -h
    vs file golangci -h%s

`
