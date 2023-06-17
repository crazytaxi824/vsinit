package typescript

import (
	"fmt"
	"log"

	"local/src/appsettings"
	"local/src/util"
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

	fmt.Printf("\n%sInstalling packages ...%s\n", util.COLOR_GREEN, util.COLOR_RESET)

	if *jstsFlags.ESlintLocal {
		if *jstsFlags.Jest {
			// 'npm install -D ts-jest @types/jest'  // 项目用 packages, 下载到项目目录.
			eslintDependencies = append(eslintDependencies, "ts-jest", "@types/jest")
		}

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

		if *jstsFlags.Jest {
			// 'npm install -D ts-jest @types/jest'  // 项目用 packages, 下载到项目目录.
			if err := util.NpmInstallDevLocal([]string{"ts-jest", "@types/jest"}); err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}
