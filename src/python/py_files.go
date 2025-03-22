package python

import (
	"fmt"

	"local/src/files"
	"local/src/util"
)

// 一定需要写的文件
func filesNeedToWrite() []util.FileContent {
	return []util.FileContent{
		{
			Filepath: ".editorconfig",
			Content:  files.Editorconfig,
		},
		{
			Filepath: ".gitignore",
			Content:  files.Gitignore,
		},
		{
			Filepath: "pyproject.toml",
			Content:  files.PyProject,
		},
		{
			Filepath: "src/tests/py_test.py",
			Content:  files.PyTest,
		},
		{
			Filepath: "src/main.py",
			Content:  files.PyMain,
		},
	}
}

func writeProjectFiles() error {
	err := util.Prompt("Python")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = util.WriteFiles(filesNeedToWrite())
	if err != nil {
		return err
	}

	fmt.Printf(jsMsg, util.COLOR_BOLD_YELLOW, util.COLOR_RESET)
	return nil
}

func writeSingleFile() error {
	fs, err := util.ChooseSingleFile(filesNeedToWrite(), "write")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = util.WriteFiles(fs)
	if err != nil {
		return err
	}

	return nil
}

func printSingleFile() error {
	fs, err := util.ChooseSingleFile(filesNeedToWrite(), "print")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("%sfile content:%s\n"+string(fs[0].Content), util.COLOR_GREEN, util.COLOR_RESET)
	return nil
}

const jsMsg = `%srun:
	1. python3 -m venv .venv
	2. source .venv/bin/activate
	3. pip3 install debugpy%s
`
