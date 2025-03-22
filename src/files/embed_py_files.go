package files

import (
	_ "embed"
)

var (
	//go:embed py_proj_files/gitignore
	PyGitignore []byte

	//go:embed py_proj_files/main.py
	PyMain []byte

	//go:embed py_proj_files/py_test.py
	PyTest []byte

	//go:embed py_proj_files/pyproject.toml
	PyProject []byte

	// common/editorconfig
)
