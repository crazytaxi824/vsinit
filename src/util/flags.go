package util

import (
	"flag"
	"fmt"
)

// for javascript & typescript use only
type JSTSFlags struct {
	FlagSet     *flag.FlagSet
	ESlintLocal *bool // set eslintrc.json locally, default globally
	Jest        *bool // test tool
}

func SetJSTSFlags() *JSTSFlags {
	var tsfs JSTSFlags
	tsfs.FlagSet = flag.NewFlagSet("ts flags", flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// eslint
	tsfs.ESlintLocal = tsfs.FlagSet.Bool("eslint-local", false,
		"set eslint config file locally\n(default: globally)")

	// alias
	f := tsfs.FlagSet.Lookup("eslint-local")
	tsfs.FlagSet.Var(f.Value, "l", fmt.Sprintf("alias to -%s", f.Name))

	// jest
	tsfs.Jest = tsfs.FlagSet.Bool("jest", false, "install typescript jest environment")
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
	rf.FlagSet = flag.NewFlagSet("ts flags", flag.ExitOnError) // Call os.Exit(2) or for -h/-help Exit(0)

	// eslint
	rf.ESlintLocal = rf.FlagSet.Bool("eslint-local", false,
		"set eslint config file locally\n(default: globally)")

	// alias
	f := rf.FlagSet.Lookup("eslint-local")
	rf.FlagSet.Var(f.Value, "l", fmt.Sprintf("alias to -%s", f.Name))

	return &rf
}
