// NOTE this command line tool only works at PROJECT ROOT directory.
// This command line tool is used to initialize project of different
// program languages in vscode.

package main

import (
	"flag"
	"fmt"
	"os"

	"local/src/golang"
	"local/src/js"
	"local/src/python"
	"local/src/ts"
	"local/src/util"
)

const languages = "go/py/ts/js/react"

func main() {
	if len(os.Args) < 2 {
		helpMsg()
		os.Exit(2)
	}

	// flag.ExitOnError will os.Exit(2) if subcommand Parse() error.
	tsCmd := flag.NewFlagSet("ts", flag.ExitOnError)
	jest := tsCmd.Bool("jest", false, "add 'jest' - unit test components")

	var folders []string
	var files []util.FileContent

	switch os.Args[1] {
	case "go":
		folders = golang.CreateFolders
		files = golang.FilesAndContent
		fmt.Println("init Golang project")
	case "py":
		folders = python.CreateFolders
		files = python.FilesAndContent
		fmt.Println("init Python project")
	case "ts":
		// parse arges first
		// nolint // flag.ExitOnError will do the os.Exit(2)
		tsCmd.Parse(os.Args[2:])

		folders = ts.CreateFolders
		files = ts.FilesAndContent
		if *jest {
			folders = append(folders, ts.JestFolder)  // add "test" folder
			files = append(files, ts.JestFileContent) // add jest example test file

			// 设置 jest
			err := ts.JestSetup()
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}

		fmt.Println("init TypeScript project")
	case "react":
		folders = ts.CreateFolders
		files = ts.ReactFilesAndContent
		fmt.Println("init React - TypeScript project")
	case "js":
		folders = js.CreateFolders
		files = js.FilesAndContent
		fmt.Println("init JavaScript project")
	default:
		helpMsg()
		os.Exit(2)
	}

	// create folders and write files
	util.WriteCfgFiles(folders, files)
}

func helpMsg() {
	fmt.Println("please specify language -", languages)
	fmt.Println("usage: vsinit <language> [<args>]")
	fmt.Println("eg: vsinit go")
}
