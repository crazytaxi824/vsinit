package ts

const (
	// 整个 golangci-lint 的设置占位符
	lintPlaceHolder = "${eslintPlaceHolder}"

	// go.lintFlags --config 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// golangci 文件夹
	golangciDirector = "/golangci"

	// golangci-lint config file path
	devciFilePath  = "/dev-ci.yml"
	prodciFilePath = "/prod-ci.yml"

	// vscode workspace
	vsWorkspace = "${workspaceRoot}"
)

// golangci-lint setting
var (
	eslintconfig = []byte(`
  // 在 OUTPUT -> ESlint 频道打印 debug 信息. 用于配置 eslint.
  "eslint.debug": true,

  // save 的时候运行 eslint
  "eslint.run": "onSave",

  // eslint 检查文件类型
  "eslint.validate": [
    "typescriptreact",
    "typescript",
    "javascriptreact",
    "javascript"
  ],

  // 单独设置 eslint 配置文件
  "eslint.options": {
    // NOTE eslint(cmd)<=v7.x 可以工作，但是 CLIEngine 已经弃用。
    // https://eslint.org/docs/developer-guide/nodejs-api#cliengine
    // 这里是全局 eslint 配置文件的固定地址
    // ts lint config file
    "configFile": "` + configPlaceHolder + `"
  },
`)
)
