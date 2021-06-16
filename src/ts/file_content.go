package ts

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": settingsJSONContent,
	".eslintrc.json":        eslintContent,
	"tsconfig.json":         tsConfigContent,
	"src/main.ts":           mainFileContent,
	".gitignore":            gitignoreContent,
}

// file content
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
  // "eslint.options": {
  //   "configFile": ".eslintrc.json"
  // },

  // NOTE important, ts string 单引号
  "prettier.singleQuote": true,

  // search.exclude 用来忽略搜索的文件夹
  // files.exclude 用来忽略工程打开的文件夹
  // 直接写文件/文件夹名字就实在项目根路径下进行匹配，不要用 / ./ 开头，
  // **/所有路径下进行匹配
  "search.exclude": {
    ".idea": true,
    "*.iml": true,
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
  "include": ["src/**/*"], // 指定编译文件
  "exclude": ["node_modules", "out", "dist", "**/*.spec.ts", "**/*.config.js"] // 排除编译文件
}  
`

	eslintContent = `/*
https://eslint.org/docs/rules/
https://www.npmjs.com/package/@typescript-eslint/eslint-plugin
https://www.npmjs.com/package/eslint-config-airbnb-typescript
https://github.com/iamturns/create-exposed-app/blob/master/.eslintrc.js
  
  需要安装的库:
  // @typescript-eslint/parser

  eslint-plugin-import
  eslint-plugin-jsx-a11y
  eslint-plugin-react
  eslint-plugin-react-hooks
  @typescript-eslint/eslint-plugin
  eslint-config-airbnb-typescript
  
  如果需要使用 react:
  "extends": "airbnb-typescript"
  如果不需要使用 react:
  "extends": "airbnb-typescript/base"
*/
{
  "extends": [
    // "eslint:recommended",
    // "airbnb-typescript", // js 也可以用, with react
    "airbnb-typescript/base", // js 也可以用
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/recommended-requiring-type-checking"
  ],
  "env": { "node": true, "browser": true, "es6": true },
  "plugins": ["@typescript-eslint"],
  // parser config
  // "parser": "@typescript-eslint/parser",
  "parserOptions": {
    "ecmaVersion": 6, // es6 = es2015
    "ecmaFeatures": {
      "jsx": true
    },
    "impliedStrict": true,
    "sourceType": "module", // script | module
    "project": "./tsconfig.json" // NOTE important, ts项目用
  },
  // 不需要检查的文件
  "ignorePatterns": [
    ".vscode",
    "out",
    "dist",
    "node_modules",
    "**/vendor/*.js",
    "**/vendor/*.ts"
  ],
  "rules": {
    "no-console": 0, // DEBUG use only.
    "no-prototype-builtins": "off",
    "import/prefer-default-export": "off",
    // "import/no-default-export": "error",
    // NOTE react use only
    // "react/destructuring-assignment": "off",
    // "react/jsx-filename-extension": "off",
    // "curly": "error", // FIXME not working, 强制 if/for/do/while 使用一致的括号风格
    "no-bitwise": "off", // 不允许使用特殊运算符 &, &=, |, |=, ^, ^=, <<, <<=, >>, >>=, >>>, >>>=, ~
    "complexity": "warn", // default 20
    // NOTE lack of no-null-keyword checks.

    // 功能检查
    "max-len": ["warn", { "code": 120 }],
    "max-lines": ["warn", 500], // 文件不超过 n 行

    // TS
    // https://github.com/typescript-eslint/typescript-eslint/blob/HEAD/packages/eslint-plugin/docs/rules/member-ordering.md
    "@typescript-eslint/member-ordering": "warn", // class 中 member 排序
    "@typescript-eslint/consistent-type-definitions": "error", // 使用统一的类型定义
    "@typescript-eslint/no-use-before-define": [
      "error",
      {
        "functions": false,
        "classes": true,
        "variables": true,
        "typedefs": true
      }
    ],
    // interface 必须 I 开头
    // https://github.com/typescript-eslint/typescript-eslint/blob/HEAD/packages/eslint-plugin/docs/rules/naming-convention.md
    "@typescript-eslint/naming-convention": [
      "warn",
      {
        "selector": "interface",
        "format": ["PascalCase"],
        "custom": {
          "regex": "^I[A-Z]",
          "match": true
        }
      }
    ]
  },
  "overrides": [
    {
      "files": ["*.js"],
      "rules": {
        // Allow require()
        "@typescript-eslint/no-var-requires": "off"
      }
    }
  ]
}
`
)
