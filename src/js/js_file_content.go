// 纯 js 项目用, 不包含 react lint

package js

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"src/main.js":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// file content
const (
	launchJSONContent = `{
  "version": "0.2.0",
  "configurations": [
    {
      // run single file
      "name": "current file",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${file}"
    },
    {
      // launch project
      "name": "src/main.js",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${workspaceFolder}/src/main.js"
    }
  ]
}
`

	settingsJSONContent = `{
  // https://eslint.org/docs/developer-guide/nodejs-api#cliengine
  // 单独设置 eslint 配置文件
  "eslint.options": {
    // 这里是全局 eslint 配置文件的固定地址
    "configFile": "/Users/ray/projects/lints/ts/eslintrc-js.json"
  },

  // eslint 检查文件类型
  "eslint.validate": [
	  // "typescriptreact",
    // "typescript",
	  // "javascriptreact",
    "javascript"
  ],

  // NOTE important, ts string 单引号
  "prettier.singleQuote": true,

  // search.exclude 用来忽略搜索的文件夹
  // files.exclude 用来忽略工程打开的文件夹
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "search.exclude": {
    ".idea": true,
    "*.iml": true,
    "out": true,
    "dist": true,
    "**/vendor": true,
    "node_modules": true,
    ".vscode": true,
    ".history": true
  },
	  
  // files.exclude 不显示文件，
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "files.exclude": {
    ".idea": true,
  }
}
`

	gitignoreContent = `# http://git-scm.com/docs/gitignore
# 项目根路径下使用 "/" 开头，如果不写 "/" 则在整个项目中进行匹配，类似 "**/"
/.vscode
/.idea
/*.iml
/.history
/node_modules

# 配置文件
/config.*

# 任何路径下用 **/ 开头
**/debug
**/vendor
`

	mainFileContent = `function main() {
  console.log("hello world");
}
  
main();
`
)
