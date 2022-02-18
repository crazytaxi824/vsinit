// this package for writing a single '.editorconfig' file only.
// all functions in 'util' copied here.

package singlefile

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"local/src/resource"
	"local/src/util"
	"log"
	"os"
	"strings"
)

var ceFlags *editorConfigFlags

// 获取命令行的前两个参数 - eg: "vs go"
func getCmd() string {
	return fmt.Sprintf("%q", strings.Join(os.Args[:2], " "))
}

// '.editorconfig' flags
type editorConfigFlags struct {
	FlagSet   *flag.FlagSet
	Overwrite *bool
}

func setEditorConfigFlags() *editorConfigFlags {
	var ecf editorConfigFlags
	ecf.FlagSet = flag.NewFlagSet(getCmd(), flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// overwrite '.editorconfig' file
	ecf.Overwrite = ecf.FlagSet.Bool("overwrite", false, "overwrite '.editorconfig' file")

	// alias
	f := ecf.FlagSet.Lookup("overwrite")
	ecf.FlagSet.Var(f.Value, "over", fmt.Sprintf("alias to -%s\n", f.Name))

	// 这里没有任何 flag, 这里只是为了 -h 命令.
	return &ecf
}

func filesNeedToWrite() []util.FileContent {
	return []util.FileContent{
		{
			FileName:  ".editorconfig",
			Content:   resource.Editorconfig,
			Overwrite: *ceFlags.Overwrite,
		},
	}
}

func writeProjectFiles() error {
	return util.WriteAllFiles(filesNeedToWrite())
}

func WriteEditorConfigFile() error {
	// go flags only for `-help`
	ceFlags = setEditorConfigFlags()
	err := ceFlags.FlagSet.Parse(os.Args[2:])
	if err != nil {
		log.Println(err)
		return err
	}

	// ask before init project
	err = askBeforeProceed(".editorconfig")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// write project files
	err = writeProjectFiles()
	if err != nil {
		return err
	}

	return nil
}

// 询问是否在当前文件夹初始化项目
func askBeforeProceed(lang string) error {
	pwd, err := os.Getwd() // 获取当前路径
	if err != nil {
		return errors.New(writeFileCanceled(lang))
	}

	fmt.Printf("Write file %s%q%s at %s%q%s? [Yes/no]: ",
		util.COLOR_BOLD_YELLOW, lang, util.COLOR_RESET, util.COLOR_YELLOW, pwd, util.COLOR_RESET)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return errors.New(writeFileCanceled(lang))
	}

	if input != "yes\n" && input != "Yes\n" {
		return errors.New(writeFileCanceled(lang))
	}

	return nil
}

// 打印红色 cancel 信息
func writeFileCanceled(lang string) string {
	return fmt.Sprintf("%sWrite file %q Canceled!%s", util.COLOR_RED, lang, util.COLOR_RESET)
}
