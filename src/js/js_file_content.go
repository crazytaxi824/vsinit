// 纯 js 项目用, 不包含 react lint

package js

import _ "embed" // for go:embed file use

var CreateFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSONContent []byte

	//go:embed cfgfiles/settings.json
	settingsJSONContent []byte

	//go:embed cfgfiles/gitignore
	gitignoreContent []byte
)

var FilesAndContent = map[string][]byte{
	".vscode/launch.json":   launchJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"src/main.js":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// file content
var mainFileContent = []byte(`main();

function main() {
  console.log('hello world');
}
`)
