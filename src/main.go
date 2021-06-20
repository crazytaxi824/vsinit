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
		fmt.Println("please specify language -", languages)
		fmt.Println("usage: vsinit <language> [<args>]")
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	jest := testCmd.Bool("jest", false, "add 'jest' - unit test components")

	var folders []string
	var files []util.FileContent

	switch os.Args[1] {
	case "go":
		fmt.Println("init Golang project")
		folders = golang.CreateFolders
		files = golang.FilesAndContent
	case "py":
		fmt.Println("init Python project")
		folders = python.CreateFolders
		files = python.FilesAndContent
	case "ts":
		fmt.Println("init TypeScript project")
		folders = ts.CreateFolders
		files = ts.FilesAndContent
	case "react":
		fmt.Println("init React - TS project")
		folders = ts.CreateFolders
		files = ts.ReactFilesAndContent
	case "js":
		fmt.Println("init JavaScript project")
		folders = js.CreateFolders
		files = js.FilesAndContent
	case "test":
		err := testCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("test command parse flag error:", err)
			return
		}
		fmt.Println("jest flag is", *jest)
		// fmt.Println("this is a command test function")
		return
	default:
		fmt.Println("languang supported -", languages)
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	// create folders and write files
	util.WriteCfgFiles(folders, files)
}
