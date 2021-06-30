package util

import (
	"os"
	"os/exec"
)

type DependenciesInstall struct {
	Dependencies []string
	Prefix       string
	Global       bool
}

// npm install libs to devDependencies
// 指定位置安装 eslint 所需依赖
func NpmInstallDependencies(path string, global bool, libs ...string) error {
	if len(libs) == 0 {
		return nil
	}

	// TODO 是否需要安装？(y/n)

	var args []string
	if global {
		args = []string{"i", "-g"}
	} else {
		args = []string{"i", "-D"}
	}

	// 指定下载到什么地方
	if path != "" {
		// --prefix 将 node_modules 创建到 path下. <path>/node_modules
		args = append(args, "--prefix", path)
	}

	// 执行命令
	args = append(args, libs...)
	cmd := exec.Command("npm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
