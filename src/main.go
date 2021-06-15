package main

import (
	"fmt"
	"os"

	"local/src/golang"
	"local/src/python"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("please specify language - go/py/ts/js")
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}

	switch args[1] {
	case "go":
		fmt.Println("init Golang project")
		golang.GoCfgFile()
	case "py":
		fmt.Println("init Python project")
		python.PyCfgFile()
	case "ts":
		fmt.Println("init TypeScript project")
	case "js":
		fmt.Println("init JavaScript project")
	default:
		fmt.Println("languang supported - go/py/ts/js")
		fmt.Println("eg: vsinit go")
		os.Exit(2)
	}
}
