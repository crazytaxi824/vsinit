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

	// for unit test 'jest' use
	//go:embed cfgfiles/example.test.js
	exampleTestJS []byte
)

// file content
var mainJS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

// FilesAndContent JS project files
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

// JestFileContent add example of unit test
var JestFileContent = util.FileContent{
	Path:    "test/example.test.js",
	Content: exampleTestJS,
}
