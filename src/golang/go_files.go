// 根据 flags 选择 files

package golang

import (
	"local/src/resource"
	"local/src/util"
)

func filesNeedToWrite() []util.FileContent {
	return []util.FileContent{
		{
			Dir:      ".vscode/",
			FileName: "settings.json",
			Content:  resource.GoVsSettings,
		},
		{
			Dir:      ".vscode/",
			FileName: "launch.json",
			Content:  resource.GoVsLaunch,
		},
		{
			FileName: ".gitignore",
			Content:  resource.GoGitignore,
		},
		{
			FileName: ".golangci.yml",
			Content:  resource.Golangci,
		},
		{
			FileName: ".editorconfig",
			Content:  resource.Editorconfig,
		},
		{
			Dir:      "src/",
			FileName: "main.go",
			Content:  resource.GoMain,
		},
	}
}

func writeProjectFiles() error {
	return util.WriteAllFiles(filesNeedToWrite())
}
