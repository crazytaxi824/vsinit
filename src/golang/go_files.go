// 根据 flags 选择 files

package golang

import (
	"fmt"

	"local/src/files"
	"local/src/util"
)

func filesNeedToWrite() []util.FileContent {
	return []util.FileContent{
		{
			Filepath: ".nvim/settings.lua",
			Content:  files.GoNvimSettings,
		},
		{
			Filepath: ".vscode/settings.json",
			Content:  files.GoVsSettings,
		},
		{
			Filepath: ".vscode/launch.json",
			Content:  files.GoVsLaunch,
		},
		{
			Filepath: ".editorconfig",
			Content:  files.Editorconfig,
		},
		{
			Filepath: ".gitignore",
			Content:  files.GoGitignore,
		},
		{
			Filepath: ".golangci.yml",
			Content:  files.Golangci,
		},
		{
			Filepath: "src/main.go",
			Content:  files.GoMain,
		},
	}
}

func writeProjectFiles() error {
	err := util.Prompt("Go")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = util.WriteFiles(filesNeedToWrite())
	if err != nil {
		return err
	}

	fmt.Printf(goMsg, util.COLOR_BOLD_YELLOW, util.COLOR_RESET)
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

const goMsg = `%srun:
    go mod init <repo>%s
`
