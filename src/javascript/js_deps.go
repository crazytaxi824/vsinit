package javascript

import (
	"fmt"
	"local/src/appsettings"
	"local/src/util"
	"log"
)

// return dependencies' names
func installPackages() error {
	// eslint dependencies
	var eslintDependencies = []string{
		"eslint-plugin-jest",        // jest unit test
		"eslint-plugin-promise",     // promise 用法
		"eslint-config-prettier",    // 解决 vscode 插件中 prettier 造成的代码问题
		"eslint-config-airbnb-base", // js 专用 lint
		// "eslint-config-airbnb",   // js 专用 lint, dependencies includes "eslint-config-airbnb-base"
	}

	fmt.Printf("\n%sInstalling packages ...%s\n", util.COLOR_GREEN, util.COLOR_RESET)

	if *jstsFlags.ESlintLocal {
		// `npm install -D <packages>`
		if err := util.NpmInstallDevLocal(eslintDependencies); err != nil {
			log.Println(err)
			return err
		}
	} else {
		// `npm install -D --prefix <global_path> <packages>`
		if err := util.NpmInstallPrefixDev(appsettings.ESLintGlobalPath, eslintDependencies); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
