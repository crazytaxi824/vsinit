package resource

import _ "embed"

var (
	//go:embed js_proj_files/gitignore
	JSGitignore []byte

	//go:embed js_proj_files/settings.json
	JSVsSettings []byte

	//go:embed js_proj_files/launch.json
	JSVsLaunch []byte

	//go:embed js_proj_files/eslintrc-js.json
	JSESlint []byte

	//go:embed js_proj_files/main.js
	JSMain []byte

	//go:embed js_proj_files/example.test.js
	JSTest []byte

	//go:embed js_proj_files/package.json
	JSPackageJSON []byte

	// common/editorconfig
)
