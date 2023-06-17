package react

import (
	"fmt"
	"log"
	"os"

	"local/src/util"
)

var reactFlags *util.ReactFlags

func InitReactProj() error {
	// parse flags
	reactFlags = util.SetReactFlags()
	err := reactFlags.FlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Println(err)
		return err
	}

	// ask before init project
	err = util.AskBeforeProceed("React")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Printf("%sIniting React Project ...%s\n", util.COLOR_GREEN, util.COLOR_RESET)

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
