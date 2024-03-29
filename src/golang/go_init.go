// install '.golangci.yml' locally for vim-go.

package golang

import (
	"fmt"
	"log"
	"os"

	"local/src/util"
)

func InitGoProj() error {
	// go flags only for `-help`
	goFlags := util.SetGoFlags()
	err := goFlags.FlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Println(err)
		return err
	}

	// ask before init project
	err = util.Prompt("Go")
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
