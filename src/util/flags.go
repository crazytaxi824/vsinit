package util

import "flag"

type TSJSFlags struct {
	FlagSet                   *flag.FlagSet
	Jest, ESLint, ESLintLocal *bool
}

func SetupTSJSFlags() TSJSFlags {
	var tsjs TSJSFlags
	tsjs.FlagSet = flag.NewFlagSet("ts/js", flag.ExitOnError)
	tsjs.Jest = tsjs.FlagSet.Bool("jest", false, "setup 'jest' in the Project")
	tsjs.ESLint = tsjs.FlagSet.Bool("eslint", false, "setup 'eslint' globally")
	tsjs.ESLintLocal = tsjs.FlagSet.Bool("eslint-local", false, "setup 'eslint' in the Project")

	return tsjs
}

type GoFlags struct {
	FlagSet             *flag.FlagSet
	Cilint, CilintLocal *bool
}

func SetupGoFlags() GoFlags {
	var gofs GoFlags
	gofs.FlagSet = flag.NewFlagSet("go", flag.ExitOnError)
	gofs.Cilint = gofs.FlagSet.Bool("cilint", false, "setup 'golangci-lint' globally")
	gofs.CilintLocal = gofs.FlagSet.Bool("cilint-local", false, "setup 'golangci-lint' in this Project")

	return gofs
}
