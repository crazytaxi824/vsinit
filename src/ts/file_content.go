package ts

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/tasks.json":    tasksJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"tslint.json":           tslintContent,
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
      "name": "current file",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${file}",
      "preLaunchTask": "tsc: build - tsconfig.json",
      "outFiles": ["${workspaceFolder}/dist/**/*.js"]
    },
    {
      "name": "src/main.ts",
      "type": "node",
      "request": "launch",
      "skipFiles": ["<node_internals>/**"],
      "program": "${workspaceFolder}/src/main.ts",
      "preLaunchTask": "tsc: build - tsconfig.json",
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
/dist
/node_modules

# 配置文件
/config.*

# 任何路径下用 **/ 开头
**/debug
**/vendor
`

	tslintContent = `{
  "extends": "tslint:recommended",
  "no-null-keyword": true, // 不允许使用null,使用undefined代替null，指代空指针对象

  // TS特性
  "rules": {
    "no-console": false, // 不允许使用console对象
    "member-access": true, // 设置成员对象的访问权限（public,private,protect)
    "member-ordering": [
      // 设置修饰符顺序
      true,
      {
        "order": [
          "static-field",
          "instance-field",
          "static-method",
          "instance-method"
        ]
      }
    ],
    "no-parameter-reassignment": true, // 不允许修改方法输入参数

    // 功能特性
    "await-promise": true, // 不允许没有Promise的情况下使用await
    "curly": true, // if/for/do/while强制使用大括号
    "no-bitwise": false, // 不允许使用特殊运算符 &, &=, |, |=, ^, ^=, <<, <<=, >>, >>=, >>>, >>>=, ~
    "no-for-in-array": true, // 不允许对Array使用for-in
    "no-invalid-template-strings": true, // 只允许在模板字符串中使用${
    "no-invalid-this": true, // 不允许在class之外使用this.
    "no-sparse-arrays": true, // 不允许array中有空元素.
    "no-switch-case-fall-through": true, // 不允许 case 在没有使用 break 的情况下 fall through 到下一个 case.
    "restrict-plus-operands": true, // 不允许自动类型转换，如果已设置不允许使用关键字var该设置无效
    "triple-equals": true, // 必须使用恒等号，进行等于比较

    // 维护性功能
    "indent": [true, "spaces", 4], // 每行开始以4个空格符开始
    "linebreak-style": [true], // 换行符格式 CR/LF可以通用使用在windows和osx
    "max-file-line-count": [true, 500], // 定义每个文件代码行数
    "max-line-length": [true, 120], // 定义每行代码数
    "no-default-export": true, // 禁止使用export default关键字，因为当export对象名称发生变化时，需要修改import中的对象名。https://github.com/palantir/tslint/issues/1182#issue-151780453
    "no-duplicate-imports": true, // 禁止在一个文件内，多次引用同一module
    "align": [
      // 定义对齐风格
      true,
      "parameters",
      "arguments",
      "statements",
      "members",
      "elements"
    ],
    "encoding": true, // 定义编码格式默认utf-8
    "interface-name": [true, "always-prefix"] // interface必须以I开头
  }
}
`
	tsConfigContent = `{
  "$schema": "https://json.schemastore.org/tsconfig",
  "compileOnSave": false,
  "compilerOptions": {
    "module": "commonjs",
    "moduleResolution": "node",
    "target": "es2015", // es6=es2015, 默认情况下使用es6，拥有 map & set
    "lib": ["es2017", "DOM"], // 包含了声明文件列表，你仍然拥有较新的类型声明
    "outDir": "./dist", // ts转为js文件时的地址
    "sourceMap": true, // 必须为true
    "strict": true, // alwaysStrict, noImplicitAny, noImplicitThis, strictBindCallApply, strictFunctionTypes, strictNullChecks and strictPropertyInitialization
    "allowJs": true,
    "noEmit": true,
    "experimentalDecorators": true, // 使用装饰器
    "importHelpers": true,
    "downlevelIteration": true,
    "allowSyntheticDefaultImports": true,
    "esModuleInterop": true,
    "skipLibCheck": false,
    // "isolatedModules": true,
    "jsx": "react-native"
  },
  "include": ["src/**/*"], // 指定编译文件
  "exclude": ["node_modules", "**/*.spec.ts", "**/*.config.js"] // 排除编译文件
}  
`

	mainFileContent = `function main() {
  console.log("hello world");
}
  
main();
`
)
