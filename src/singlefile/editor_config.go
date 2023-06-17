// this package for writing a single '.editorconfig' file only.
// all functions in 'util' copied here.

package singlefile

import (
	"flag"
	"fmt"
	"log"
	"os"

	"local/src/resource"
	"local/src/util"
)

var ecFlags *editorConfigFlags

// '.editorconfig' flags
type editorConfigFlags struct {
	FlagSet   *flag.FlagSet
	Overwrite *bool
}

func setEditorConfigFlags() *editorConfigFlags {
	var ecf editorConfigFlags
	ecf.FlagSet = flag.NewFlagSet(getCmd(), flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// overwrite '.editorconfig' file
	ecf.Overwrite = ecf.FlagSet.Bool("overwrite", false, "overwrite '.editorconfig' file\n")

	// alias
	f := ecf.FlagSet.Lookup("overwrite")
	ecf.FlagSet.Var(f.Value, "ov", fmt.Sprintf("alias to -%s\n", f.Name))

	return &ecf
}

func editorConfigFile() []util.FileContent {
	return []util.FileContent{
		{
			FileName:  ".editorconfig",
			Content:   resource.Editorconfig,
			Overwrite: *ecFlags.Overwrite,
		},
	}
}

func writeEditorConfigFiles() error {
	return util.WriteAllFiles(editorConfigFile())
}

func WriteEditorConfigFile() error {
	ecFlags = setEditorConfigFlags()
	err := ecFlags.FlagSet.Parse(os.Args[3:])
	if err != nil {
		log.Println(err)
		return err
	}

	// ask before write file
	err = askBeforeProceed(".editorconfig")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// write files
	err = writeEditorConfigFiles()
	if err != nil {
		return err
	}

	return nil
}
