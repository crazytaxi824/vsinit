package ts

import _ "embed" // for go:embed file use

var (
	//go:embed tscfgfiles/launch.json
	launchJSONContent string

	//go:embed tscfgfiles/settings.json
	settingsJSONContent string

	//go:embed tscfgfiles/tasks.json
	tasksJSONContent string

	//go:embed tscfgfiles/gitignore
	gitignoreContent string

	//go:embed tscfgfiles/tsconfig.json
	tsConfigContent string

	//go:embed reactcfgfiles/settings.json
	reactSettingsJSONContent string
)

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

var ReactFilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": reactSettingsJSONContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// ts file content
const (
	mainFileContent = `function main() {
  console.log('hello world');
}

main();
`
)
