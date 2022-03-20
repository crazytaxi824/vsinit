package singlefile

import (
	"flag"
	"fmt"
	"local/src/resource"
	"local/src/util"
	"log"
	"os"
)

var gciFlags *golangciFlags

// '.golangci.yml' flags
type golangciFlags struct {
	FlagSet   *flag.FlagSet
	Overwrite *bool
	NoGeneric *bool // false (default) - settings for go1.18; true - settings for go1.17
}

func setGolangciFlags() *golangciFlags {
	var gcif golangciFlags
	gcif.FlagSet = flag.NewFlagSet(getCmd(), flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// overwrite '.golangci.yml' file
	gcif.Overwrite = gcif.FlagSet.Bool("overwrite", false, "overwrite '.golangci.yml' file\n")

	// no-generic flag
	gcif.NoGeneric = gcif.FlagSet.Bool("no-generic", false,
		"'golangci.yml' no-generic settings, for go version < 1.18\n")

	// alias
	overwrite := gcif.FlagSet.Lookup("overwrite")
	gcif.FlagSet.Var(overwrite.Value, "ov", fmt.Sprintf("alias to -%s\n", overwrite.Name))

	nogeneric := gcif.FlagSet.Lookup("no-generic")
	gcif.FlagSet.Var(nogeneric.Value, "ng", fmt.Sprintf("alias to -%s\n", nogeneric.Name))

	return &gcif
}

func golangciFile() []util.FileContent {
	// 判断是否需要使用 no-generic settings
	if *gciFlags.NoGeneric {
		return []util.FileContent{
			{
				FileName:  ".golangci.yml",
				Content:   resource.Golangci17,
				Overwrite: *gciFlags.Overwrite,
			},
		}
	}

	return []util.FileContent{
		{
			FileName:  ".golangci.yml",
			Content:   resource.Golangci,
			Overwrite: *gciFlags.Overwrite,
		},
	}
}

func writeGolangciFiles() error {
	return util.WriteAllFiles(golangciFile())
}

func WriteGolangciFile() error {
	gciFlags = setGolangciFlags()
	err := gciFlags.FlagSet.Parse(os.Args[3:])
	if err != nil {
		log.Println(err)
		return err
	}

	// ask before write file
	err = askBeforeProceed(".golangci.yml")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// write file
	err = writeGolangciFiles()
	if err != nil {
		return err
	}

	return nil
}
