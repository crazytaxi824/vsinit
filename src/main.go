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

const languages = "go/py/ts/js"

func main() {
	if len(os.Args) < 2 {
		helpMsg()
		os.Exit(2)
	}

	// flag.ExitOnError will os.Exit(2) if subcommand Parse() error.
	tsjsCmd := flag.NewFlagSet("tsjs", flag.ExitOnError)
	jestflag := tsjsCmd.Bool("jest", false, "add 'jest' - unit test components")

	switch os.Args[1] {
	case "go":
		folders := golang.CreateFolders
		files := golang.FilesAndContent

		fmt.Println("init Golang project")
		util.WriteCfgFiles(folders, files)
	case "py":
		folders := python.CreateFolders
		files := python.FilesAndContent

		fmt.Println("init Python project")
		util.WriteCfgFiles(folders, files)
	case "ts":
		// parse arges first
		// nolint // flag.ExitOnError will do the os.Exit(2)
		tsjsCmd.Parse(os.Args[2:])

		folders := ts.CreateFolders
		files := ts.FilesAndContent

		var npmLibs []string // Dependencies needs to be downloaded

		if *jestflag {
			// add jest example test file
			folders = append(folders, ts.TestFolder)
			files = append(files, ts.JestFileContent)

			// 设置 jest
			var err error
			npmLibs, err = ts.SetupTS()
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}
		}

		// write project files first
		fmt.Println("init TypeScript project")
		util.WriteCfgFiles(folders, files)

		// then npm install after wirte package.json file
		if err := util.NpmInstallDependencies(npmLibs...); err != nil {
			os.Exit(2)
		}

	case "js":
		// parse arges first
		// nolint // flag.ExitOnError will do the os.Exit(2)
		tsjsCmd.Parse(os.Args[2:])

		folders := js.CreateFolders
		files := js.FilesAndContent
		if *jestflag {
			// add jest example test file
			folders = append(folders, js.TestFolder)
			files = append(files, js.JestFileContent)
		}

		fmt.Println("init JavaScript project")
		util.WriteCfgFiles(folders, files)
	default:
		helpMsg()
		os.Exit(2)
	}
}

func helpMsg() {
	fmt.Println("please specify language -", languages)
	fmt.Println("usage: vsinit <language> [<args>]")
	fmt.Println("eg: vsinit go")
}
