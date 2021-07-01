package util

import "flag"

type TSJSFlags struct {
	FlagSet                   *flag.FlagSet
	Jest, ESLint, ESLintLocal *bool
}

func SetupTSJSFlags() TSJSFlags {
	var tsjs TSJSFlags
	tsjs.FlagSet = flag.NewFlagSet("ts/js", flag.ExitOnError)
	tsjs.Jest = tsjs.FlagSet.Bool("jest", false, "add 'jest' locally")
	tsjs.ESLint = tsjs.FlagSet.Bool("eslint", false, "add 'eslint' globally")
	tsjs.ESLintLocal = tsjs.FlagSet.Bool("eslint-local", false, "add 'eslint' in the Project")

	return tsjs
}

type GoFlags struct {
	FlagSet             *flag.FlagSet
	Cilint, CilintLocal *bool
}

func SetupGoFlags() GoFlags {
	var gofs GoFlags
	gofs.FlagSet = flag.NewFlagSet("go", flag.ExitOnError)
	gofs.Cilint = gofs.FlagSet.Bool("cilint", false, "add 'golangci-lint' globally")
	gofs.CilintLocal = gofs.FlagSet.Bool("cilint-local", false, "add 'golangci-lint' in this Project")

	return gofs
}
