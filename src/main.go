// NOTE this command line tool only works at PROJECT ROOT directory.
// This command line tool is used to initialize project of different
// program languages in vscode.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"local/src/golang"
	"local/src/js"
	"local/src/python"
	"local/src/ts"
	"local/src/util"
)

const languages = "go/py/ts/js"

func helpMsg() {
	fmt.Println("please specify language -", languages)
	fmt.Println("usage: vsinit <language> [<args>]")
	fmt.Println("eg: vsinit go")
}

func main() {
	if len(os.Args) < 2 {
		helpMsg()
		os.Exit(2)
	}

	// flag.ExitOnError will os.Exit(2) if subcommand Parse() error.
	tsjsSet := flag.NewFlagSet("ts/js", flag.ExitOnError)
	jestflag := tsjsSet.Bool("jest", false, "add 'jest' - unit test components")
	eslintflag := tsjsSet.Bool("eslint", false, "setup eslint globally")
	eslintProjectflag := tsjsSet.Bool("eslint-proj", false, "setup eslint in project")

	goSet := flag.NewFlagSet("go", flag.ExitOnError)
	cilintflag := goSet.Bool("cilint", false, "setup golangci-lint globally")
	cilintProjectflag := goSet.Bool("cilint-proj", false, "setup golangci-lint in project")

	var (
		err         error
		suggestions []*util.Suggestion
	)
	switch os.Args[1] {
	case "go":
		suggestions, err = golang.InitProject(goSet, cilintflag, cilintProjectflag)
	case "py":
		err = python.InitProject()
	case "ts":
		suggestions, err = ts.InitProject(tsjsSet, jestflag, eslintflag, eslintProjectflag)
	case "js":
		err = js.InitProject(tsjsSet, jestflag)
	case "envcheck":
		suggestions, err = golang.CheckGO(false) // DEBUG
	default:
		helpMsg()
		os.Exit(2)
	}

	// 统一打印 error
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// 打印提醒
	printSuggestions(suggestions)
}

func printSuggestions(suggestions []*util.Suggestion) {
	if len(suggestions) == 0 {
		fmt.Println("All Done! Happy Hunting.")
		return
	}

	var builder strings.Builder
	for _, sug := range suggestions {
		builder.WriteString(sug.String())
	}
	fmt.Print(builder.String()) // 这里用 println 会多空一行
}
