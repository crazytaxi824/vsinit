package golang

import (
	_ "embed" // for go:embed file use
	"local/src/util"
)

// setup golangci-lint.yml config file.
func SetupLocalGolangciLint(ciDir string) (folders []string, files []util.FileContent) {
	folders = []string{ciDir, ciDir + util.GolangciDirector}
	files = []util.FileContent{
		{Path: ciDir + util.GolangciDirector + "/dev-ci.yml", Content: devci},
		{Path: ciDir + util.GolangciDirector + "/prod-ci.yml", Content: prodci},
	}
	return
}
