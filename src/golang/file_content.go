package golang

// file content
const (
	launchJSONContent = `{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Auto Main",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "port": 12345,
      "host": "127.0.0.1",
      "program": "${workspaceRoot}/src", // main.go 路径
      "cwd": "${workspaceRoot}",		 // 只在 debug 模式时有用
      // "env": {},
      // "args": ["-c","/xxx/config.yml"],
      "internalConsoleOptions": "openOnSessionStart",
      "showLog": true // show logs in debug mode
    }
  ]
}
`

	settingsJSONContent = `{
  // golangci-lint 单独设置
  // "go.lintFlags": ["--config=~/.golangci/release-ci.yml"],
  // "go.lintOnSave": "package",

  // search.exclude 用来忽略搜索的文件夹
  // files.exclude 用来忽略工程打开的文件夹
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "search.exclude": {
    ".idea": true,
    ".vscode": true,
    "*.iml": true,
    "**/vendor": true,
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

# 配置文件
/config.*

# 任何路径下用 **/ 开头
**/*.gorun
**/debug
**/vendor
**/go.sum
`

	mainFileContent = `package main
`
)
