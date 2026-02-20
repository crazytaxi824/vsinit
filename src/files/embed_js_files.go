package files

import (
	_ "embed"
)

var (
	//go:embed js_proj_files/settings.json
	JSVsSettings []byte

	//go:embed js_proj_files/launch.json
	JSVsLaunch []byte

	//go:embed js_proj_files/main.js
	JSMain []byte

	//go:embed js_proj_files/example.test.js
	JSTest []byte

	//go:embed js_proj_files/jsconfig.json
	JsConfig []byte
)
