package javascript

import (
	"log"
	"os"

	"local/src/util"
)

func InitProj() error {
	flags := util.SetFlags()
	err := flags.FlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Println(err)
		return err
	}

	// choose file to print
	if *flags.Print {
		err = printSingleFile()
		if err != nil {
			return err
		}
		return nil
	}

	// choose file to write
	if *flags.File {
		err = writeSingleFile()
		if err != nil {
			return err
		}
		return nil
	}

	// write all files
	err = writeProjectFiles()
	if err != nil {
		return err
	}

	return nil
}
