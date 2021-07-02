package util

import (
	"os"
	"os/exec"
)

// 需要安装的依赖
type dependenciesInstall struct {
	dependencies []string
	prefix       string
	global       bool
}

// npm install libs to devDependencies
// 指定位置安装 eslint 所需依赖
func npmInstallDependencies(di dependenciesInstall) error {
	if len(di.dependencies) == 0 {
		return nil
	}

	var args []string
	if di.global {
		args = []string{"i", "-g"}
	} else {
		args = []string{"i", "-D"}
	}

	// 指定下载到什么地方
	if di.prefix != "" {
		// --prefix 将 node_modules 创建到 path下. <path>/node_modules
		args = append(args, "--prefix", di.prefix)
	}

	// 执行命令
	args = append(args, di.dependencies...)
	cmd := exec.Command("npm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
