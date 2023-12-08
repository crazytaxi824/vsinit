package util

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// 获取命令行的前两个参数 - eg: "vs go"
func getCmd() string {
	return fmt.Sprintf("%q", strings.Join(os.Args[:2], " "))
}

type Flags struct {
	FlagSet *flag.FlagSet

	File  *bool // choose a single file to write
	Print *bool // choose a single file to print
}

func SetFlags() *Flags {
	var flags Flags
	flags.FlagSet = flag.NewFlagSet(getCmd(), flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)
	flags.File = flags.FlagSet.Bool("file", false, "choose single file to write")
	flags.Print = flags.FlagSet.Bool("print", false, "choose single file to print")

	// alias
	f := flags.FlagSet.Lookup("file")
	flags.FlagSet.Var(f.Value, "f", fmt.Sprintf("alias to -%s\n", f.Name))

	p := flags.FlagSet.Lookup("print")
	flags.FlagSet.Var(p.Value, "p", fmt.Sprintf("alias to -%s\n", p.Name))

	// golang 没有任何 flag, 这里只是为了 -h 命令.
	return &flags
}
