// 纯 js 项目用, 不包含 react lint

package js

import _ "embed" // for go:embed file use

var (
	//go:embed cfgfiles/launch.json
	launchJSONContent string

	//go:embed cfgfiles/settings.json
	settingsJSONContent string

	//go:embed cfgfiles/gitignore
	gitignoreContent string
)

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"src/main.js":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// file content
const (
	mainFileContent = `main();

function main() {
  console.log('hello world');
}
`
)
