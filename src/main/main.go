package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	//获得当前路径
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 检查vscode and src 文件夹是否存在
	f, err := ioutil.ReadDir(curPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var vsExist = false
	var gitExist = false
	var srcExist = false

	for _, v := range f {
		if !v.IsDir() {
			if v.Name() == ".gitignore" {
				gitExist = true
			}
		} else {
			if v.Name() == ".vscode" {
				vsExist = true
			} else if v.Name() == "src" {
				srcExist = true
			}
		}
	}

	if !vsExist {
		err = os.Mkdir(".vscode", 0777)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if !srcExist {
		err = os.MkdirAll("src/main", 0777)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	// 分别判断 launch, settings, tasks 是否存在
	var launchExist = false
	var settingsExist = false
	var tasksExist = false

	vs, err := ioutil.ReadDir(curPath + "/.vscode")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range vs {
		if v.Name() == "launch.json" {
			launchExist = true
		} else if v.Name() == "settings.json" {
			settingsExist = true
		} else if v.Name() == "tasks.json" {
			tasksExist = true
		}
	}

	// 写文件
	if !launchExist {
		launchPath := curPath + "/.vscode/launch.json"
		writeLaunch(launchPath)
	}

	if !settingsExist {
		settingsPath := curPath + "/.vscode/settings.json"
		writeSettings(settingsPath)
	}

	if !tasksExist {
		tasksPath := curPath + "/.vscode/tasks.json"
		writeTasks(tasksPath)
	}

	if !srcExist {
		mainPath := curPath + "/src/main/main.go"
		writeMain(mainPath)
	}

	if !gitExist {
		ignorePath := curPath + "/.gitignore"
		writeIgnore(ignorePath)
	}

}

func writeLaunch(launchPath string) {
	// launch.json
	launch := `{
	// 使用 IntelliSense 了解相关属性。 
	// 悬停以查看现有属性的描述。
	// 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
	"version": "0.2.0",
	"configurations": [
		{
			"name": "Launch",
			"type": "go",
			"request": "launch",
			"mode": "exec",
			"remotePath": "",
			"port": 2345,
			"host": "127.0.0.1",
			"program": "${workspaceRoot}/debug",
			"preLaunchTask": "go build",
			"internalConsoleOptions": "openOnSessionStart",
			"env": {
				"GOPATH": "${workspaceRoot}"
			},
			"args": [],
			"showLog": true
		}
	]
}`

	err := ioutil.WriteFile(launchPath, []byte(launch), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("/.vscode/launch.json 写入完成")
}

func writeSettings(settingsPath string) {
	settings := `{
	//search.exclude 用来忽略搜索的文件夹
	//files.exclude 用来忽略工程打开的文件夹
	"files.exclude": {
		"${workspaceRootFolderName}_.gorun": true,
		"${workspaceRootFolderName}_bin": true,
		"${workspaceRootFolderName}.iml": true,
		"debug":true,
	},
	//设置gopath
	"go.gopath": "${workspaceRoot}",
	//"go.toolsGopath":”/Users/ray/go”
}`
	err := ioutil.WriteFile(settingsPath, []byte(settings), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("/.vscode/settings.json 写入完成")
}

func writeTasks(tasksPath string) {
	tasks := `{
    // See https://go.microsoft.com/fwlink/?LinkId=733558
    // for the documentation about the tasks.json format
    "version": "2.0.0",
    "tasks": [
        {
            "label": "go build",
            "type": "shell",
            "presentation": {
                "echo": true,
                "reveal": "never",
                "focus": false,
                "panel": "shared"
            },
            "command": "export",
            "args": [
                "GOPATH=${workspaceRoot};",
                "go",
                "build",
                "-o",
                "debug",
                "./src/main/"
            ]
        }
    ]
}`
	err := ioutil.WriteFile(tasksPath, []byte(tasks), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("/.vscode/tasks.json 写入完成")
}

func writeIgnore(ignorePath string) {
	ignore := `/.vscode
/.idea
/*.iml
/pkg
/*.gorun
/debug`

	err := ioutil.WriteFile(ignorePath, []byte(ignore), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("/.gitignore 写入完成")
}

func writeMain(mainPath string) {
	mainFile := `package main

func main() {
	
}`

	err := ioutil.WriteFile(mainPath, []byte(mainFile), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("/src/main/main.go 写入完成")
}
