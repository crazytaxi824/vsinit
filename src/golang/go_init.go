// copy files: .vscode/settints.json, .vscode/launch.json, .editorconfig, .gitignore, .golangci.yml, src/main.go, src/main_test.go
// install '.golangci.yml' locally for vim-go.

package golang

import (
	"fmt"
	"local/src/util"
)

func InitGoProj() error {
	// ask before init project
	err := util.AskBeforeProceed("Go")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("%sIniting Go Project ...%s\n", util.COLOR_GREEN, util.COLOR_RESET)

	// write project files
	err = writeProjectFiles()
	if err != nil {
		return err
	}

	return nil
}
