package python

// file content
const (
	launchJSONContent = `{
  "version": "0.2.0",
  "configurations": [
    {
      // use for single file
      "name": "current file",
      "type": "python",
      "request": "launch",
      "program": "${file}",
      "console": "integratedTerminal"
    },
    {
      // use for project
      "name": "src/main.py",
      "type": "python",
      "request": "launch",
      "program": "${workspaceRoot}/src/main.py",
      "console": "integratedTerminal"
    }
  ]
}
`

	settingsJSONContent = `{
  // 选择 python version
  "python.pythonPath": "/usr/local/bin/python3",

  // search.exclude 用来忽略搜索的文件夹
  // files.exclude 用来忽略工程打开的文件夹
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "search.exclude": {
    ".idea": true,
    "*.iml": true,
    "**/vendor": true,
    ".vscode": true,
    ".mypy_cache": true,
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
**/debug
**/vendor
`

	mainGoContent = `def main():
  print("hello world")


main()
`
)
