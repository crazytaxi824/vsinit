// copy files: .vscode/settints.json, .vscode/launch.json, .vim/coc-settings.json,
//    .editorconfig, .gitignore, .golangci.yml, src/main.go, src/main_test.go
// install '.golangci.yml' locally for vim-go.

package resource

import (
	_ "embed"
)

var (
	//go:embed go_proj_files/gitignore
	GoGitignore []byte

	//go:embed go_proj_files/settings.json
	GoVsSettings []byte

	//go:embed go_proj_files/launch.json
	GoVsLaunch []byte

	//go:embed go_proj_files/golangci.yml
	Golangci []byte

	//go:embed go_proj_files/golangci-1.17.yml
	Golangci17 []byte

	//go:embed go_proj_files/main_file.go.txt
	GoMain []byte

	// common/editorconfig
)
