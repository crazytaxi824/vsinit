package ts

import (
	_ "embed" // for go:embed file use

	"local/src/util"
)

var CreateFolders = []string{".vscode", "src"}

var (
	//go:embed tscfgfiles/launch.json
	launchJSONContent []byte

	//go:embed tscfgfiles/settings.json
	settingsJSONContent []byte

	//go:embed tscfgfiles/tasks.json
	tasksJSONContent []byte

	//go:embed tscfgfiles/gitignore
	gitignoreContent []byte

	//go:embed tscfgfiles/tsconfig.json
	tsConfigContent []byte

	//go:embed reactcfgfiles/settings.json
	reactSettingsJSONContent []byte
)

var mainFileContent = []byte(`main();

function main() {
  console.log('hello world');
}
`)

var FilesAndContent = []util.FileContent{
	{
		Path:    ".vscode/launch.json",
		Content: launchJSONContent,
	},
	{
		Path:    ".vscode/tasks.json",
		Content: tasksJSONContent,
	},
	{
		Path:    ".vscode/settings.json",
		Content: settingsJSONContent,
	},
	{
		Path:    ".gitignore",
		Content: gitignoreContent,
	},
	{
		Path:    "tsconfig.json",
		Content: tsConfigContent,
	},
	{
		Path:    "src/main.ts",
		Content: mainFileContent,
	},
}

var ReactFilesAndContent = []util.FileContent{
	{
		Path:    ".vscode/launch.json",
		Content: launchJSONContent,
	},
	{
		Path:    ".vscode/tasks.json",
		Content: tasksJSONContent,
	},
	{
		// 主要修改是 setting，里面改变了 lint 的 config 文件地址
		Path:    ".vscode/settings.json",
		Content: reactSettingsJSONContent,
	},
	{
		Path:    ".gitignore",
		Content: gitignoreContent,
	},
	{
		Path:    "tsconfig.json",
		Content: tsConfigContent,
	},
	{
		Path:    "src/main.ts",
		Content: mainFileContent,
	},
}
