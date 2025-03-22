package files

import (
	_ "embed"
)

//go:embed common/editorconfig
var Editorconfig []byte

//go:embed common/gitignore
var Gitignore []byte
