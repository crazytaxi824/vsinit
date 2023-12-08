package javascript

import (
	"fmt"

	"local/src/files"
	"local/src/util"
)

// 一定需要写的文件
func filesNeedToWrite() []util.FileContent {
	return []util.FileContent{
		{
			Filepath: ".nvim/settings.lua",
			Content:  files.JSNvimSettings,
		},
		{
			Filepath: ".vscode/settings.json",
			Content:  files.JSVsSettings,
		},
		{
			Filepath: ".vscode/launch.json",
			Content:  files.JSVsLaunch,
		},
		{
			Filepath: ".editorconfig",
			Content:  files.Editorconfig,
		},
		{
			Filepath: ".gitignore",
			Content:  files.JSGitignore,
		},
		{
			Filepath: ".eslintrc.json",
			Content:  files.JSESlint,
		},
		{
			Filepath: "package.json",
			Content:  files.JSPackageJSON, // jest settings included
		},
		{
			Filepath: "example.test.js",
			Content:  files.JSTest,
		},
		{
			Filepath: "src/main.js",
			Content:  files.JSMain,
		},
	}
}

func writeProjectFiles() error {
	err := util.Prompt("Javascript")
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
    npm install -g eslint jest
    npm install -D <packages> # eslint deps: eslint-config-prettier
".eslintrc.json" file:
    change "settings.jest.version" to your current jest version.%s
`
