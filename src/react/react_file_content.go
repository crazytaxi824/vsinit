// NOTE 先使用 'npx create-react-app xxx --template typescript' 创建项目，然后使用 'vs init react'.
// 修改 package.json - eslintConfig.extends 中添加 "./eslintrc-react.json" 地址
// 修改 tsconfig.json - compilerOptions 中添加 "noImplicitReturns": true, // 强制定义出参类型
// 修改 .gitignore 文件，添加 /.vscode  /eslintrc-react.json
// 添加 .vscode/settings.json, eslintrc-react.json 文件, 不需要 launch.json, 启动 react 项目运行 `npm run start`
// npm 安装 eslint 依赖

package react

import (
	"bytes"
	_ "embed" // for go:embed file use
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"local/src/util"
	"os"
)

// eslint dependencies
var eslintDependencies = []string{
	"@typescript-eslint/eslint-plugin", // 必须
	"eslint-config-airbnb-typescript",  // ts 用
	"eslint-plugin-jest",               // jest unit test
	"eslint-plugin-promise",            // promise 用法
	"eslint-config-prettier",           // 解决 vscode 插件中 prettier 造成的代码问题
	"@types/react-router-dom",          // react-router-dom typescript
}

var ErrCreateReactFirst = errors.New("need to create-react-app first")

var (
	//go:embed cfgfiles/settings.json
	settingsJSON []byte

	//go:embed cfgfiles/eslintrc-react.json
	eslintrcJSON []byte

	//go:embed cfgfiles/tsconfig.json
	tsconfigJSON []byte

	// ESLint 配置文件名
	eslintFilePath = "/eslintrc-react.json" // NOTE JS 要改
)

var (
	createFolders = []string{".vscode"}

	filesAndContent = []util.FileContent{
		{Path: util.SettingsJSONPath, Content: settingsJSON},
		{Path: "eslintrc-react.json", Content: eslintrcJSON},
	}
)

func InitProject() ([]*util.Suggestion, error) {
	// command args 检查
	if len(os.Args) > 3 && (os.Args[3] == "-h" || os.Args[3] == "-help" || os.Args[3] == "--help") {
		fmt.Println("Usage of 'vs init react': no flags")
		os.Exit(0)
	} else if len(os.Args) > 3 {
		util.HelpMsg()
		os.Exit(2)
	}

	// 添加 .vscode/settings.json, eslintrc-react.json 文件
	ff := util.InitFoldersAndFiles(createFolders, filesAndContent)

	// 修改 gitignore 文件
	err := changeGitignore(ff)
	if err != nil {
		return nil, err
	}

	// 修改 tsconfig.json 文件
	err = checkTsconfig(ff)
	if err != nil {
		return nil, err
	}

	// 修改 package.json 文件
	packageSuggestion(ff)

	// 下载 dependencies
	err = ff.AddMissingDependencies(eslintDependencies, "package.json", "")
	if err != nil {
		return nil, err
	}

	// 写入所需文件
	fmt.Println("init TypeScript project")
	if err := ff.WriteAllFiles(); err != nil {
		return nil, err
	}

	// 安装所有缺失的依赖
	if err := ff.InstallMissingDependencies(); err != nil {
		return nil, err
	}

	// 返回 suggestion
	return ff.Suggestions(), nil
}

// .gitignore 添加 /.vscode & /eslintrc-react.json
func changeGitignore(ff *util.FoldersAndFiles) error {
	gf, err := os.Open(util.GitignorePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		return ErrCreateReactFirst
	}

	r, err := io.ReadAll(gf)
	if err != nil {
		return err
	}

	// 内容添加在最后
	r = append(r, []byte("\n# dev environment\n/.vscode\n"+eslintFilePath+"\n")...)

	ff.AddFiles(util.FileContent{
		Path:      util.GitignorePath,
		Content:   r,
		Overwrite: true,
	})

	return nil
}

type tsconfig struct {
	CompilerOptions struct {
		Target                           string   `json:"target,omitempty"`
		Lib                              []string `json:"lib,omitempty"`
		AllowJS                          bool     `json:"allowJs,omitempty"`
		SkipLibCheck                     bool     `json:"skipLibCheck,omitempty"`
		EsModuleInterop                  bool     `json:"esModuleInterop,omitempty"`
		AllowSyntheticDefaultImports     bool     `json:"allowSyntheticDefaultImports,omitempty"`
		Strict                           bool     `json:"strict,omitempty"`
		ForceConsistentCasingInFileNames bool     `json:"forceConsistentCasingInFileNames,omitempty"`
		NoFallthroughCasesInSwitch       bool     `json:"noFallthroughCasesInSwitch,omitempty"`
		Module                           string   `json:"module,omitempty"`
		ModuleResolution                 string   `json:"moduleResolution,omitempty"`
		ResolveJSONModule                bool     `json:"resolveJsonModule,omitempty"`
		IsolatedModules                  bool     `json:"isolatedModules,omitempty"`
		NoEmit                           bool     `json:"noEmit,omitempty"`
		Jsx                              string   `json:"jsx,omitempty"`
	} `json:"compilerOptions,omitempty"`
	Include []string `json:"include,omitempty"`
}

// 主要目的是为了检查 create-react-app 新版本中对 tsconfig.json 是否有调整。
func checkTsconfig(ff *util.FoldersAndFiles) error {
	tf, err := os.Open("tsconfig.json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		return ErrCreateReactFirst
	}

	jsonc, err := io.ReadAll(tf)
	if err != nil {
		return err
	}

	// tsconfig.json 转 json
	j, err := util.JSONCToJSON(jsonc)
	if err != nil {
		return err
	}

	// tsconfig.json 文件存在，通过 json 反序列化检查设置和之前版本有没有变化
	var tsc tsconfig
	de := json.NewDecoder(bytes.NewReader(j))
	de.DisallowUnknownFields() // 如果 json 中有字段而结构体中没有则会报错
	err = de.Decode(&tsc)
	if err != nil {
		ff.AddSuggestions(&util.Suggestion{
			Problem: fmt.Sprintf("error: %s, please add following to 'tsconfig.json' file", err.Error()),
			Solution: `"noImplicitReturns": true, // 强制定义出参类型
			"exclude": ["node_modules", "**/*.spec.ts", "**/*.config.js"]`,
		})
		return nil
	}

	// 如果没有问题，则直接写入新的 tsconfig.json 文件
	ff.AddFiles(util.FileContent{
		Path:      "tsconfig.json",
		Content:   tsconfigJSON,
		Overwrite: true,
	})

	return nil
}

// 建议修改 package.json 文件，package.json 里面可能会有很多东西需要调整，就不用程序调整了。
func packageSuggestion(ff *util.FoldersAndFiles) {
	ff.AddSuggestions(&util.Suggestion{
		Problem: "please add following to 'package.json' eslintConfig:",
		Solution: `  "eslintConfig": {
    "extends": [
      "react-app",
      "react-app/jest"
    ],
    "rules": {
      "@typescript-eslint/no-unused-vars": "off"
    }
  },`,
	})
}
