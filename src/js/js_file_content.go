// 纯 js 项目用, 不包含 react lint

package js

import (
	_ "embed" // for go:embed file use

	"local/src/util"
)

var CreateFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSONContent []byte

	//go:embed cfgfiles/settings.json
	settingsJSONContent []byte

	//go:embed cfgfiles/gitignore
	gitignoreContent []byte
)

// file content
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
		Path:    ".vscode/settings.json",
		Content: settingsJSONContent,
	},
	{
		Path:    ".gitignore",
		Content: gitignoreContent,
	},
	{
		Path:    "src/main.js",
		Content: mainFileContent,
	},
}
