package main

import (
	"fmt"
	"local/src/golang"
	"local/src/javascript"
	"local/src/react"
	"local/src/singlefile"
	"local/src/typescript"
	"local/src/util"
	"os"
)

func main() {
	// 只输入了 'vs' 命令的情况下
	if len(os.Args) < 2 {
		fmt.Printf(helpMsg, util.COLOR_YELLOW, util.COLOR_RESET, util.COLOR_YELLOW, util.COLOR_RESET)
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

	case "editorconfig", "editor":
		if err := singlefile.WriteEditorConfigFile(); err != nil {
			os.Exit(2)
		}

	default:
		fmt.Printf(helpMsg, util.COLOR_YELLOW, util.COLOR_RESET, util.COLOR_YELLOW, util.COLOR_RESET)
		os.Exit(2)
	}

	fmt.Printf("%sAll Done! Happy Coding!%s\n", util.COLOR_BOLD_GREEN, util.COLOR_RESET)
}

const helpMsg = `help:%s
    vs [go|ts(typescript)|js(javascript)|react|editor(editorconfig)]%s

more info:%s
    vs ts -h%s

`
