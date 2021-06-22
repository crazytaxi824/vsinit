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

	//go:embed cfgfiles/package.json
	packageJSON []byte

	//go:embed cfgfiles/tsconfig.json
	tsconfigJSON []byte

	//go:embed cfgfiles/example.test.ts
	exampleTestTS []byte
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
		Path:    "package.json",
		Content: packageJSON,
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
