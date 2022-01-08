// react - eslint flag

package react

import (
	"bytes"
	"fmt"
	"local/src/appsettings"
	"local/src/resource"
	"local/src/util"
	"log"
	"strings"
)

// files need to write
var fs []util.FileContent

// 一定需要写的文件
func filesHasToWrite() {
	fs = append(fs, []util.FileContent{
		{
			FileName:   ".gitignore",
			Content:    resource.ReactGitignore,
			Suggestion: fmt.Sprintf("# add dev environment\n%s/.vscode\n**/*.bak%s\n", util.COLOR_YELLOW, util.COLOR_RESET), // add config
		},
		{
			FileName:   "tsconfig.json",
			Content:    resource.ReactConfigJSON,
			Suggestion: fmt.Sprintf(tsconfigSuggestion, util.COLOR_YELLOW, util.COLOR_RESET, util.COLOR_YELLOW, util.COLOR_RESET),
		},
		{
			FileName: ".editorconfig",
			Content:  resource.Editorconfig,
		},
		{
			// package.json 中需要修改设置，不需要重新写入文件.
			FileName:   "package.json",
			Content:    []byte("{}"),
			Suggestion: fmt.Sprintf(packageSuggestion, util.COLOR_YELLOW, util.COLOR_RESET, util.COLOR_YELLOW, util.COLOR_RESET),
		},
	}...)
}

// vscode 'settings.json' content 'eslint config filepath' based on ESLint flag - local / global
func filesMightNeedToWrite() {
	var overrideConfigFile string

	if *reactFlags.ESlintLocal {
		// vscode settings -> "overrideConfigFile": "本地相对位置"
		overrideConfigFile = fmt.Sprintf(`"overrideConfigFile": %q`, appsettings.ReactESLintFileName)

		// eslint 文件安装在项目本地
		fs = append(fs,
			util.FileContent{
				// Dir:      "", // dir path based on ESLint flag - local / global
				FileName: appsettings.ReactESLintFileName,
				Content:  resource.ReactESlint,
			},
		)
	} else {
		// vscode settings -> "overrideConfigFile": "全局绝对位置"
		overrideConfigFile = fmt.Sprintf(`"overrideConfigFile": %q`, appsettings.ESLintGlobalPath+appsettings.ReactESLintFileName)

		// eslint 文件安装在 Global Path 位置
		fs = append(fs,
			util.FileContent{
				Dir:      appsettings.ESLintGlobalPath, // dir path based on ESLint flag - local / global
				FileName: appsettings.ReactESLintFileName,
				Content:  resource.ReactESlint,
			},
		)
	}

	// 添加 .vscode/settings.json 文件
	fs = append(fs,
		util.FileContent{
			Dir:      ".vscode/",
			FileName: "settings.json",
			Content: bytes.ReplaceAll(resource.ReactVsSettings,
				[]byte(`"overrideConfigFile": "eslintrc-react.json"`), // 这里是写死在 .vscode/settings.json 文件中的内容, 不要改.
				[]byte(overrideConfigFile),
			),
			Suggestion: fmt.Sprintf(settingsSuggestion, util.COLOR_YELLOW, overrideConfigFile, util.COLOR_RESET),
		},
	)
}

const (
	// package.json
	packageSuggestion = `{
  // 注意: package.json 不能写注释
  %s"proxy": "http://localhost:18080",%s  // 本地代理, 绕过 CORS
  "eslintConfig": {
    "extends": ["react-app", "react-app/jest"],
    %s"rules": {
      "@typescript-eslint/no-unused-vars": "off"
    }%s
  },
}
`
	// tsconfig.json
	tsconfigSuggestion = `{
  "compilerOptions": {
    %s"noImplicitReturns": true,%s // 强制定义出参类型
  },
  // 排除检查&编译文件
  %s"exclude": ["node_modules", "**/*.spec.ts", "**/*.config.js"]%s
}
`
	// .vscode/settings.json eslint filepath suggestion
	settingsSuggestion = `{
  "eslint.options": {%s
    %s%s
  },
}
`
)

// react common function files
func commonFiles() ([]util.FileContent, error) {
	var cfs []util.FileContent

	fsDir := "react_proj_files/react_common_fn"

	fsd, err := resource.ReactCommonFuncs.ReadDir(fsDir)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, f := range fsd {
		if strings.Contains(f.Name(), ".tsx") {
			content, er := resource.ReactCommonFuncs.ReadFile(fmt.Sprintf("%s/%s", fsDir, f.Name()))
			if er != nil {
				log.Println(er)
				return nil, er
			}

			cfs = append(cfs, util.FileContent{
				Dir:      "src/util/",
				FileName: f.Name(),
				Content:  content,
			})
		}
	}

	return cfs, nil
}

func writeProjectFiles() error {
	filesHasToWrite()
	filesMightNeedToWrite()

	cfs, err := commonFiles()
	if err != nil {
		return err
	}

	err = util.WriteAllFiles(fs)
	if err != nil {
		return err
	}

	return util.WriteAllFiles(cfs)
}
