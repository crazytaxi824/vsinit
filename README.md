# 初始化 go / js / ts / python 项目

## go

- `.vscode/settings.json`

- `.vscode/launch.json`

- `.gitignore`

- `.golangci.yml`

  - 全局地址固定位置 `~/.config/vsi/golangci.yml` // 可以解析 `~/` 路径
  - 本地地址固定位置 `${workspaceRoot}/.golangci.yml` // 可以解析 `${workspaceRoot}` 路径, 这是 vscode 提供的环境变量.

- `src/main.go`

- `.editorconfig` for github tab_size

<br />

## js - javascript

- `.vscode/settings.json`

  - 全局地址固定位置 `/Users/ray/.config/vsi/eslintrc-js.json` // 绝对地址必须 `/` 开头, 否则会被认为是相对 `项目根目录` 的地址. eg: `~/xxx` 会被解析成 `project/~/xxx` 地址.
  - 本地地址固定位置 `eslintrc-js.yml` // 相对地址是项目根目录. 无法解析 `${workspaceRoot}` 路径.
  - eslint 工作原理: 首先需要 eslint 命令行工具 `npm install -g eslint`, eslint 会按照上述指定路径查找 `package.json` 文件, 最后按照 `package.json` 查找相关规则. 所以 `package.json` 必须存在. 需要指定 `--prefix <path>` 安装 eslint rules, `npm install -D --prefix <path> <packages>`

- `.vscode/launch.json`

- `.gitignore`

- `eslint.json` - `npm install eslint rules`

- `package.json`

- `src/main.js`

- `src/example.test.js` - jest

- `.editorconfig` for github tab_size

<br />

## ts - typescript

- `.vscode/settings.json`

- `.vscode/launch.json`

- `.vscode/tasks.json`

- `tsconfig.json` - 比 js 多了 tsconfig 编译设置. 指定了将 ts 编译成 js 的规则.

- `.gitignore`

- `eslint.json` - 同 `js` 设置, 只是需要安装的 eslint rules 不同

- `package.json`

- `src/main.js`

- `src/example.test.ts` - jest

- `.editorconfig` for github tab_size

- `jest` 设置 - `npm install -D ts-jest @types/jest`; js 不需要安装

<br />

## react - ts

react - CRA (create react app) 项目中包含了其他文件.

- `.vscode/settings.json` - 不需要 launch

- `tsconfig.json`

- `eslint.json` - `npm install -D <eslint rules>`

- 需要修改 `package.json` 中的设置让自定义 eslint 生效. `package.json` 文件不能写 comments.

主要是添加 "rules", 防止 auto compile 的时候因为 unused rule 而报错.

```json
{
  ...
  "eslintConfig": {
    "extends": [
      ...
    ],
    "rules": {
      "@typescript-eslint/no-unused-vars": "off"
    }
  }
  ...
}
```

- 修改 `tsconfig.json` 文件.

```json
{
  ...
  "compilerOptions": {
    "noImplicitReturns": true, // 强制定义出参类型
  },
  // 排除检查&编译文件
  "exclude": ["node_modules", "**/*.spec.ts", "**/*.config.js"]
}
```

- `.editorconfig` for github tab_size

<br />

## python ? - 测试用 (Beta)

- `.vscode/settings.json`

- `.vscode/launch.json`

- `.gitignore`

<br />

# TODO

react flags:

- common-packages

- common-functions



