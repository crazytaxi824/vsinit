package javascript

import (
	"fmt"
	"local/src/util"
	"log"
	"os"
)

// TODO change name
var jstsFlags *util.JSTSFlags

func InitJSProj() error {
	// parse flags
	jstsFlags = util.SetJSTSFlags()
	err := jstsFlags.FlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Println(err)
		return err
	}

	// ask before init project
	err = util.AskBeforeProceed("Javascript")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("%sIniting Javascript Project ...%s\n", util.COLOR_GREEN, util.COLOR_RESET)

	// choose files based on flags
	err = writeProjectFiles()
	if err != nil {
		return err
	}

	// [VVI] need to write files before npm install <packages>, because 'package.json' file might be overwritten.
	// choose dependencies based on flags
	err = installPackages()
	if err != nil {
		return err
	}

	util.PrintSuggestions()

	return nil
}
