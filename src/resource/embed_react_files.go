package resource

import (
	"embed"
)

var (
	//go:embed react_proj_files/gitignore
	ReactGitignore []byte

	//go:embed react_proj_files/settings.json
	ReactVsSettings []byte

	//go:embed react_proj_files/eslintrc-react.json
	ReactESlint []byte

	//go:embed react_proj_files/tsconfig.json
	ReactConfigJSON []byte

	// common/editorconfig

	//go:embed react_proj_files/react_common_fn/*
	ReactCommonFuncs embed.FS
)
