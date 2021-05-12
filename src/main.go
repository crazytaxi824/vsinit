// 创建 vscode/launch.json, vscode/settings.json, src/main.go, .gitignore 文件
package main

import (
	"errors"
	"fmt"
	"os"
)

var allFiles = [4]string{".vscode/launch.json", ".vscode/settings.json", "src/main.go", ".gitignore"}

func main() {
	// create .vscode & src Dir
	fmt.Printf("creating .vscode & src directories ... ")
	err := createDirs()
	if err != nil {
		fmt.Println("fail")
		fmt.Println(err)
		return
	}
	fmt.Println("done")

	for i := range allFiles {
		err = createFiles(allFiles[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// create .vscode & src dir,
func createDirs() error {
	err := os.Mkdir(".vscode", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("create .vscode Dir error: %w", err)
	}

	err = os.Mkdir("src", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("create src Dir error: %w", err)
	}

	return nil
}

// creating files
func createFiles(file string) error {
	// nolint:gosec // vsPath is checked
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create %s Files error: %w", file, err)
	}
	defer func() {
		if er := f.Close(); er != nil {
			fmt.Println(er)
			return
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("get %s File status error: %w", file, err)
	}

	// file is not empty
	if fi.Size() != 0 {
		return nil
	}

	fmt.Printf("writing file: %s ... ", file)
	// write file content
	err = writeFile(f)
	if err != nil {
		fmt.Println("fail")
		return fmt.Errorf("write file %s error: %w", file, err)
	}
	fmt.Println("done")

	return nil
}

// write file content
func writeFile(file *os.File) error {
	switch file.Name() {
	case ".vscode/launch.json":
		_, err := file.Write([]byte(launchJSONContent))
		return err
	case ".vscode/settings.json":
		_, err := file.Write([]byte(settingsJSONContent))
		return err
	case ".gitignore":
		_, err := file.Write([]byte(gitignoreContent))
		return err
	case "src/main.go":
		_, err := file.Write([]byte(mainGoContent))
		return err
	}
	return nil
}

// file content
const (
	launchJSONContent = `{
  // plugins: gopls, cweill/gotests, ramya-rao-a/go-outline
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
    // "**/pkg": true,
    "*.iml": true,
    "**/vendor": true,
    ".history": true
  },
	  
  // files.exclude 不显示文件，
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "files.exclude": {
    ".idea": true,
    // "**/pkg": true,
    "*.iml": true,
    ".history": true
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

	mainGoContent = `package main
`
)
