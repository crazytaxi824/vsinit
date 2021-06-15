package main

import (
	"fmt"
	"os"

	"local/src/golang"
	"local/src/python"
	"local/src/ts"
	"local/src/util"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("please specify language - go/py/ts")
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	var fc map[string]string

	switch args[1] {
	case "go":
		fmt.Println("init Golang project")
		fc = golang.FilesAndContent
	case "py":
		fmt.Println("init Python project")
		fc = python.FilesAndContent
	case "ts":
		fmt.Println("init TypeScript project")
		fc = ts.FilesAndContent
	// case "js":
	// 	fmt.Println("init JavaScript project")
	default:
		fmt.Println("languang supported - go/py/ts")
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	if err := util.WriteCfgFiles(fc); err != nil {
		fmt.Println(err)
	}
}
