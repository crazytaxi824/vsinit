package python

import (
	"fmt"
	"local/src/util"
)

var filesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"src/main.py":           mainFileContent,
	".gitignore":            gitignoreContent,
}

func PyCfgFile() {
	// create .vscode & src Dir
	fmt.Printf("creating .vscode & src directories ... ")
	err := util.CreateVsCodeDirs()
	if err != nil {
		fmt.Println("fail")
		fmt.Println(err)
		return
	}
	fmt.Println("done")

	for fp, fc := range filesAndContent {
		err = util.CreateAndWriteFiles(fp, fc)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
