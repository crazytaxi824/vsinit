package main

import (
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
	args := os.Args
	if len(args) < 2 {
		fmt.Println("please specify language -", languages)
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	var folders []string
	var fc map[string][]byte

	switch args[1] {
	case "go":
		fmt.Println("init Golang project")
		folders = golang.CreateFolders
		fc = golang.FilesAndContent
	case "py":
		fmt.Println("init Python project")
		folders = python.CreateFolders
		fc = python.FilesAndContent
	case "ts":
		fmt.Println("init TypeScript project")
		folders = ts.CreateFolders
		fc = ts.FilesAndContent
	case "react":
		fmt.Println("init React - TS project")
		folders = ts.CreateFolders
		fc = ts.ReactFilesAndContent
	case "js":
		fmt.Println("init JavaScript project")
		folders = js.CreateFolders
		fc = js.FilesAndContent
	default:
		fmt.Println("languang supported -", languages)
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	util.WriteCfgFiles(folders, fc)
}
