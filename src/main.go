package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"local/src/golang"
	"local/src/javascript"
	"local/src/util"
)

func main() {
	log.SetFlags(log.Llongfile)

	// 只输入了 'vs' 命令的情况下
	if len(os.Args) < 2 {
		fmt.Print(mainHelpMsg)
		os.Exit(2)
	}

	switch strings.ToLower(os.Args[1]) {
	case "go", "golang":
		if err := golang.InitProj(); err != nil {
			os.Exit(2)
		}

	case "js", "javascript":
		if err := javascript.InitProj(); err != nil {
			os.Exit(2)
		}

	default:
		fmt.Print(mainHelpMsg)
		os.Exit(2)
	}

	fmt.Printf("%sAll Done! Happy Coding!%s\n", util.COLOR_GREEN, util.COLOR_RESET)
}

const mainHelpMsg = `Init a project at current directory
Usage:
    vs [go | js]
flags info:
    vs go -h
    vs js -h
`
