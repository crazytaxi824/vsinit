package ts

import (
	_ "embed" // for go:embed file use

	"local/src/util"
)

var CreateFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSON []byte

	//go:embed cfgfiles/settings.json
	settingsJSON []byte

	//go:embed cfgfiles/tasks.json
	tasksJSON []byte

	//go:embed cfgfiles/gitignore
	gitignore []byte

	//go:embed cfgfiles/tsconfig.json
	tsconfigJSON []byte

	//go:embed reactcfgfiles/settings.json
	reactSettingsJSON []byte
)

var mainTS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

var FilesAndContent = []util.FileContent{
	{
		Path:    ".vscode/launch.json",
		Content: launchJSON,
	},
	{
		Path:    ".vscode/tasks.json",
		Content: tasksJSON,
	},
	{
		Path:    ".vscode/settings.json",
		Content: settingsJSON,
	},
	{
		Path:    ".gitignore",
		Content: gitignore,
	},
	{
		Path:    "tsconfig.json",
		Content: tsconfigJSON,
	},
	{
		Path:    "src/main.ts",
		Content: mainTS,
	},
}

var ReactFilesAndContent = []util.FileContent{
	{
		Path:    ".vscode/launch.json",
		Content: launchJSON,
	},
	{
		Path:    ".vscode/tasks.json",
		Content: tasksJSON,
	},
	{
		// 主要修改是 setting，里面改变了 lint 的 config 文件地址
		Path:    ".vscode/settings.json",
		Content: reactSettingsJSON,
	},
	{
		Path:    ".gitignore",
		Content: gitignore,
	},
	{
		Path:    "tsconfig.json",
		Content: tsconfigJSON,
	},
	{
		Path:    "src/main.ts",
		Content: mainTS,
	},
}
