// 纯 js 项目用, 不包含 react lint

package js

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

	//go:embed cfgfiles/gitignore
	gitignore []byte

	//go:embed cfgfiles/package.json
	packageJSON []byte
)

// file content
var mainJS = []byte(`main();

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
		Path:    ".vscode/settings.json",
		Content: settingsJSON,
	},
	{
		Path:    "package.json",
		Content: packageJSON,
	},
	{
		Path:    ".gitignore",
		Content: gitignore,
	},
	{
		Path:    "src/main.js",
		Content: mainJS,
	},
}
