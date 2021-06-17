package ts

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

var ReactFilesAndContent = map[string]string{
	".vscode/launch.json":   reactlaunchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": reactSettingsJSONContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// ts file content
const (
	launchJSONContent = `{
  "version": "0.2.0",
  "configurations": [
    {
      // run single file
      "name": "current file",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${file}",
      "preLaunchTask": "tsc: build - tsconfig.json",
      // "console": "integratedTerminal",
      "outFiles": ["${workspaceFolder}/dist/**/*.js"]
    },
    {
      // launch project
      "name": "src/main.ts",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${workspaceFolder}/src/main.ts",
      "preLaunchTask": "tsc: build - tsconfig.json",
      // "console": "integratedTerminal",
      "outFiles": ["${workspaceFolder}/dist/**/*.js"]
    }
  ]
}
`
	tasksJSONContent = `{
  "version": "2.0.0",
  "tasks": [
    {
      "type": "typescript",
      "tsconfig": "tsconfig.json",
      "problemMatcher": ["$tsc"], // 使用 tsc 命令转码
      "group": "build",
      "label": "tsc: build - tsconfig.json",
      "presentation": {
        "reveal": "silent",
        "clear": true,
        "showReuseMessage": false
      }
    }
  ]
}
`

	settingsJSONContent = `{
  // 单独设置 eslint 配置文件
  "eslint.options": {
    // NOTE eslint<=v7.x 可以工作，但是 CLIEngine 已经弃用。
    // 这里是全局 eslint 配置文件的固定地址
    "configFile": "/Users/ray/projects/lints/ts/eslintrc-ts.json"
  },

  // eslint 检查文件类型
  "eslint.validate": [
    "typescriptreact",  
    "typescript",
    "javascriptreact",
    "javascript"
  ],

  // NOTE important, ts string 单引号
  "prettier.singleQuote": true,

  // search.exclude 用来忽略搜索的文件夹
  // files.exclude 用来忽略工程打开的文件夹
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "search.exclude": {
    ".idea": true,
    "*.iml": true,
    "out": true,
    "dist": true,
    "**/vendor": true,
    "node_modules": true,
    ".vscode": true,
    ".history": true
  },
	  
  // files.exclude 不显示文件，
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "files.exclude": {
    ".idea": true,
  }
}
`

	gitignoreContent = `# http://git-scm.com/docs/gitignore
# 项目根路径下使用 "/" 开头，如果不写 "/" 则在整个项目中进行匹配，类似 "**/"
/.vscode
/.idea
/*.iml
/.history
/out
/dist
/node_modules

# 配置文件
/config.*

# 任何路径下用 **/ 开头
**/debug
**/vendor
`

	mainFileContent = `function main() {
  console.log("hello world");
}
  
main();
`

	tsConfigContent = `{
  // "$schema": "https://json.schemastore.org/tsconfig",
  "compileOnSave": false,
  "compilerOptions": {
    "module": "commonjs",
    "moduleResolution": "node",
    "target": "es2015", // es6=es2015, 默认情况下使用es6，拥有 map & set
    "lib": ["es2017", "DOM"], // 包含了声明文件列表，你仍然拥有较新的类型声明
    "sourceMap": true, // 必须为true，debug 用。
    "outDir": "dist", // NOTE important, ts转为js文件时的地址
    "jsx": "react", // NOTE 指定 jsx 代码的生成: preserve | react | react-native
    "allowJs": true, // 允许编译 js 文件.
    "strict": true, // 启用所有严格模式
    "noFallthroughCasesInSwitch": true, // switch 语句 fall-through 不允许
    "experimentalDecorators": true, // 使用装饰器
    "emitDecoratorMetadata": true, // 为装饰器提供元数据的支持
    "importHelpers": true,
    "downlevelIteration": true,
    "allowSyntheticDefaultImports": true,
    "esModuleInterop": true,
    "skipLibCheck": true // 跳过默认库检查
    // NOTE product 环境下需要打开以下检查.
    // "checkJs": true, // 报告 js 文件中的错误，和 allowJs 一起使用.
    // "noUnusedLocals": true, // 有未使用的局部变量时报错
    // "noUnusedParameters": true // 函数有未使用的参数时报错
    // NOTE 以下特定开发情况下再打开以下检查
    // "noEmit": true, // 只做 type check，不进行 compilation
    // "isolatedModules": true, // 开发 module, 所有 func & type 必须 import/export
  },
  // 指定检查&编译文件
  "include": ["src/**/*"], 
  // 排除检查的文件
  "exclude": ["node_modules", "out", "dist", "**/*.spec.ts", "**/*.config.js"] 
}  
`
)

// react settings
const (
	reactSettingsJSONContent = `{
  // 单独设置 eslint 配置文件
  "eslint.options": {
    // NOTE eslint<=v7.x 可以工作，但是 CLIEngine 已经弃用。
    // 这里是全局 eslint 配置文件的固定地址
    "configFile": "/Users/ray/projects/lints/ts/eslintrc-react.json"
  },

  // eslint 检查文件类型
  "eslint.validate": [
    "typescriptreact",  
    "typescript",
    "javascriptreact",
    "javascript"
  ],

  // NOTE important, ts string 单引号
  "prettier.singleQuote": true,

  // search.exclude 用来忽略搜索的文件夹
  // files.exclude 用来忽略工程打开的文件夹
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "search.exclude": {
    ".idea": true,
    "*.iml": true,
    "out": true,
    "dist": true,
    "**/vendor": true,
    "node_modules": true,
    ".vscode": true,
    ".history": true
  },
	  
  // files.exclude 不显示文件，
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "files.exclude": {
    ".idea": true,
  }
}
`

	reactlaunchJSONContent = `{
  "version": "0.2.0",
  "configurations": [
    {
      // launch project
      "name": "src/main.ts",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${workspaceFolder}/src/main.ts",
      "preLaunchTask": "tsc: build - tsconfig.json",
      // "console": "integratedTerminal",
      "outFiles": ["${workspaceFolder}/dist/**/*.js"]
    }
  ]
}
`
)
