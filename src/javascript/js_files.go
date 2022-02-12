package javascript

import (
	"bytes"
	"fmt"
	"local/src/appsettings"
	"local/src/resource"
	"local/src/util"
)

// files need to write
var fs []util.FileContent

// 一定需要写的文件
func filesHasToWrite() {
	fs = append(fs, []util.FileContent{
		{
			Dir:      ".vscode/",
			FileName: "launch.json",
			Content:  resource.JSVsLaunch,
		},
		{
			FileName: ".gitignore",
			Content:  resource.JSGitignore,
		},
		{
			FileName: ".editorconfig",
			Content:  resource.Editorconfig,
		},
		{
			Dir:      "src/",
			FileName: "main.js",
			Content:  resource.JSMain,
		},
		{
			FileName: "package.json",
			Content:  resource.JSPackageJSON, // jest settings included anyway
		},
	}...)
}

// vscode 'settings.json' content 'eslint config filepath' based on ESLint flag - local / global
// ESLint Dir path based on ESLint flag - local / global
// 'src/example.test.js' based on Jest flag
// 'package.json' content based on Jest flag
func filesMightNeedToWrite() {
	var overrideConfigFile string

	if *jstsFlags.ESlintLocal {
		// vscode settings -> "overrideConfigFile": "本地相对位置"
		overrideConfigFile = fmt.Sprintf(`"overrideConfigFile": %q`, appsettings.JSESLintFileName)

		// eslint 文件安装在项目本地
		fs = append(fs,
			util.FileContent{
				// Dir:      "", // dir path based on ESLint flag - local / global
				FileName: appsettings.JSESLintFileName,
				Content:  resource.JSESlint,
			},
		)
	} else {
		// vscode settings -> "overrideConfigFile": "全局绝对位置"
		overrideConfigFile = fmt.Sprintf(`"overrideConfigFile": %q`,
			appsettings.ESLintGlobalPath+appsettings.JSESLintFileName)

		// eslint 文件安装在 Global Path 位置
		fs = append(fs,
			util.FileContent{
				Dir:      appsettings.ESLintGlobalPath, // dir path based on ESLint flag - local / global
				FileName: appsettings.JSESLintFileName,
				Content:  resource.JSESlint,
			},
		)
	}

	// 添加 .vscode/settings.json | .vim/coc-settings.json 文件
	fs = append(fs,
		util.FileContent{
			Dir:      ".vscode/",
			FileName: "settings.json",
			Content: bytes.ReplaceAll(resource.JSVsSettings,
				[]byte(`"overrideConfigFile": "eslintrc-js.json"`), // 这里是写死在 .vscode/settings.json 文件中的内容, 不要改.
				[]byte(overrideConfigFile),
			),
			Suggestion: fmt.Sprintf(settingsSuggestion, util.COLOR_YELLOW, overrideConfigFile, util.COLOR_RESET),
		},
		util.FileContent{
			Dir:      ".vim/",
			FileName: "coc-settings.json",
			Content: bytes.ReplaceAll(resource.JSVimCocSettings,
				[]byte(`"overrideConfigFile": "eslintrc-js.json"`), // 这里是写死在 .vim/coc-settings.json 文件中的内容, 不要改.
				[]byte(overrideConfigFile),
			),
			Suggestion: fmt.Sprintf(settingsSuggestion, util.COLOR_YELLOW, overrideConfigFile, util.COLOR_RESET),
		},
	)

	if *jstsFlags.Jest {
		fs = append(fs,
			// write example.test.js file at "src/"
			util.FileContent{
				Dir:      "src/",
				FileName: "example.test.js",
				Content:  resource.JSTest,
			},
		)
	}
}

func writeProjectFiles() error {
	filesHasToWrite()
	filesMightNeedToWrite()
	return util.WriteAllFiles(fs)
}

const (
	// .vscode/settings.json eslint filepath suggestion
	settingsSuggestion = `{
  "eslint.options": {%s
    %s%s
  },
}
`
)
