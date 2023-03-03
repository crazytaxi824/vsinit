package singlefile

import (
	"flag"
	"fmt"
	"log"
	"os"

	"local/src/resource"
	"local/src/util"
)

var gciFlags *golangciFlags

// '.golangci.yml' flags
type golangciFlags struct {
	FlagSet   *flag.FlagSet
	Overwrite *bool
}

func setGolangciFlags() *golangciFlags {
	var gcif golangciFlags
	gcif.FlagSet = flag.NewFlagSet(getCmd(), flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// overwrite '.golangci.yml' file
	gcif.Overwrite = gcif.FlagSet.Bool("overwrite", false, "overwrite '.golangci.yml' file\n")

	// alias
	overwrite := gcif.FlagSet.Lookup("overwrite")
	gcif.FlagSet.Var(overwrite.Value, "ov", fmt.Sprintf("alias to -%s\n", overwrite.Name))

	return &gcif
}

func golangciFile() (string, []util.FileContent) {
	filename := ".golangci.yml"
	return filename, []util.FileContent{
		{
			FileName:  filename,
			Content:   resource.Golangci,
			Overwrite: *gciFlags.Overwrite,
		},
	}
}

func writeGolangciFiles(fc []util.FileContent) error {
	return util.WriteAllFiles(fc)
}

func WriteGolangciFile() error {
	gciFlags = setGolangciFlags()
	err := gciFlags.FlagSet.Parse(os.Args[3:])
	if err != nil {
		log.Println(err)
		return err
	}

	filename, fileContent := golangciFile()

	// ask before write file
	err = askBeforeProceed(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// write file
	err = writeGolangciFiles(fileContent)
	if err != nil {
		return err
	}

	return nil
}
