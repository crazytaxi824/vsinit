// NOTE this command line tool only works at PROJECT ROOT directory.
// This command line tool is used to initialize project of different
// program languages in vscode.

package main

import (
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
	if len(os.Args) < 3 {
		helpMsg()
		os.Exit(2)
	}

	// 接受 flags
	gofs := util.SetupGoFlags()
	tsjs := util.SetupTSJSFlags()

	var (
		err         error
		suggestions []*util.Suggestion
	)

	switch os.Args[1] {
	case "init":
		suggestions, err = initCommand(gofs, tsjs)
	case "envcheck":
		suggestions, err = envCheckCommand()
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

func initCommand(gofs util.GoFlags, tsjs util.TSJSFlags) (suggestions []*util.Suggestion, err error) {
	switch os.Args[2] {
	case "go":
		suggestions, err = golang.InitProject(gofs)
	case "py":
		err = python.InitProject()
	case "ts":
		suggestions, err = ts.InitProject(tsjs)
	case "js":
		suggestions, err = js.InitProject(tsjs)
	default:
		helpMsg()
		os.Exit(2)
	}

	if err != nil {
		return nil, err
	}

	if len(suggestions) > 0 {
		return suggestions, nil
	}
	return nil, nil
}

func envCheckCommand() (suggestions []*util.Suggestion, err error) {
	switch os.Args[2] {
	case "go":
		suggestions, err = golang.CheckGO()
	case "py":
		suggestions, err = python.CheckPython()
	case "ts":
		suggestions, err = ts.CheckTS()
	case "js":
		suggestions, err = js.CheckJS()
	default:
		helpMsg()
		os.Exit(2)
	}

	if err != nil {
		return nil, err
	}

	if len(suggestions) > 0 {
		return suggestions, nil
	}
	return nil, nil
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
