package python

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
	"src/main.py":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// file content
const (
	mainFileContent = `def main():
    print("hello world")


main()
`
)
