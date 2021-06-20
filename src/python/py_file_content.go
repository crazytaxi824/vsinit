package python

import _ "embed" // for go:embed file use

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
	"src/main.py":           mainFileContent,
	".gitignore":            gitignoreContent,
}

var mainFileContent = []byte(`def main():
    print("hello world")


main()
`)
