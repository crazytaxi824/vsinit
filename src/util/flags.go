package util

import (
	"flag"
	"fmt"
)

// golang flags.
type GoFlags struct {
	FlagSet *flag.FlagSet
}

func SetGoFlags() *GoFlags {
	var gf GoFlags
	gf.FlagSet = flag.NewFlagSet("'go' flags", flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// golang 没有任何 flag, 这里只是为了 -h 命令.
	return &gf
}

// for javascript & typescript use only
type JSTSFlags struct {
	FlagSet     *flag.FlagSet
	ESlintLocal *bool // set eslintrc.json locally, default globally
	Jest        *bool // test tool
}

func SetJSTSFlags() *JSTSFlags {
	var tsfs JSTSFlags
	tsfs.FlagSet = flag.NewFlagSet("'js/ts' flags", flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// eslint
	tsfs.ESlintLocal = tsfs.FlagSet.Bool("eslint-local", false,
		"install 'eslint-rules' related dependencies locally.\n(default: install dependencies globally)")

	// alias
	f := tsfs.FlagSet.Lookup("eslint-local")
	tsfs.FlagSet.Var(f.Value, "l", fmt.Sprintf("alias to -%s", f.Name))

	// jest
	tsfs.Jest = tsfs.FlagSet.Bool("jest", false, "install 'jest' related dependencies.")
	j := tsfs.FlagSet.Lookup("jest")
	tsfs.FlagSet.Var(j.Value, "j", fmt.Sprintf("alias to -%s", j.Name))

	return &tsfs
}

// for react use only
type ReactFlags struct {
	FlagSet     *flag.FlagSet
	ESlintLocal *bool // set eslintrc.json locally, default globally
}

func SetReactFlags() *ReactFlags {
	var rf ReactFlags
	rf.FlagSet = flag.NewFlagSet("'react' flags", flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// eslint
	rf.ESlintLocal = rf.FlagSet.Bool("eslint-local", false,
		"install 'eslint-rules' related dependencies locally.\n(default: install dependencies globally)")

	// alias
	f := rf.FlagSet.Lookup("eslint-local")
	rf.FlagSet.Var(f.Value, "l", fmt.Sprintf("alias to -%s", f.Name))

	return &rf
}
