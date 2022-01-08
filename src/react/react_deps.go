package react

import (
	"fmt"
	"local/src/appsettings"
	"local/src/util"
	"log"
)

// return dependencies' names
func installPackages() error {
	var eslintDependencies = []string{
		"@typescript-eslint/eslint-plugin", // 必须
		"@typescript-eslint/parser",        // 必须
		"eslint-config-airbnb",             // 依赖
		"eslint-config-airbnb-typescript",  // ts 用
		"eslint-plugin-jest",               // jest unit test
		"eslint-plugin-promise",            // promise 用法
		"eslint-config-prettier",           // 解决 vscode 插件中 prettier 造成的代码问题
	}

	var commonPackages = []string{
		// material-UI
		"@mui/material",
		"@mui/x-data-grid",
		"@emotion/styled",

		// router-dom
		"react-router-dom",
		"@types/react-router-dom",
		"history@4.10.1", // 指定版本号, history v4 -> react-router v4 & v5; history v5 -> react-router v6

		"react-hook-form", // 表单
		"axios",           // http 请求
		"swiper",          // 图片滚动 banner
	}

	fmt.Printf("\n%sInstalling packages ...%s\n", util.COLOR_GREEN, util.COLOR_RESET)

	if *reactFlags.ESlintLocal {
		// 合并下载 `npm install -D <packages>`
		if err := util.NpmInstallDevLocal(eslintDependencies); err != nil {
			log.Println(err)
			return err
		}
	} else {
		// 不能合并下载
		// `npm install -D --prefix <global_path> <packages>`
		if err := util.NpmInstallPrefixDev(appsettings.ESLintGlobalPath, eslintDependencies); err != nil {
			log.Println(err)
			return err
		}
	}

	// install common packages
	if err := util.NpmInstallSaveLocal(commonPackages); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
