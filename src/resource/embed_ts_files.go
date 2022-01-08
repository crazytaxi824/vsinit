package resource

import (
	_ "embed"
)

var (
	//go:embed ts_proj_files/gitignore
	TSGitignore []byte

	//go:embed ts_proj_files/settings.json
	TSVsSettings []byte

	//go:embed ts_proj_files/launch.json
	TSVsLaunch []byte

	//go:embed ts_proj_files/tasks.json
	TSVsTasks []byte

	//go:embed ts_proj_files/eslintrc-ts.json
	TSESlint []byte

	//go:embed ts_proj_files/main.ts
	TSMain []byte

	//go:embed ts_proj_files/example.test.ts
	TSTest []byte

	//go:embed ts_proj_files/package.json
	TSPackageJSON []byte

	//go:embed ts_proj_files/tsconfig.json
	TSConfigJSON []byte

	// common/editorconfig
)
