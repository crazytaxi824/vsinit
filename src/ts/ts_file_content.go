package ts

import _ "embed" // for go:embed file use

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

var FilesAndContent = map[string][]byte{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

var ReactFilesAndContent = map[string][]byte{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": reactSettingsJSONContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

var mainFileContent = []byte(`main();

function main() {
  console.log('hello world');
}
`)
