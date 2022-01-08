// args - go

package main

import (
	"fmt"
	"local/src/golang"
	"local/src/javascript"
	"local/src/react"
	"local/src/typescript"
	"local/src/util"
	"os"
)

func main() {
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

	default:
		fmt.Printf(helpMsg, util.COLOR_YELLOW, util.COLOR_RESET, util.COLOR_YELLOW, util.COLOR_RESET)
		os.Exit(2)
	}

	fmt.Printf("%sAll Done! Happy Coding!%s\n", util.COLOR_BOLD_GREEN, util.COLOR_RESET)
}

const helpMsg = `help:%s
    vs [go|ts|js|react]%s

more info:%s
    vs ts -h%s

`
