# 初始化 go / js / ts / python 项目

命令行工具

```bash
vs [go|ts|js|react]

vs <lang> -h
```

## go

flags: - 没有 flag.

files:

- `.vscode/settings.json`

  - 为了 vim-go 全部安装在本地, 不提供全局安装 `.golangci.yml`, 如果想要全局安装需要手动修改.
  - 全局地址固定位置 `~/.config/vsi/golangci.yml` // 可以解析 `~/` 路径
  - 本地地址固定位置 `${workspaceRoot}/.golangci.yml` // 可以解析 `${workspaceRoot}` 路径, 这是 vscode 提供的环境变量.

- `.vscode/launch.json`

- `.gitignore`

- `.golangci.yml`

- `src/main.go`

- `.editorconfig` for github tab_size

<br />

## js - javascript

flags:
- `--eslint-local` - 将 eslint 安装在本地.
- `--jest` - 安装 jest.

files:
- `.vscode/settings.json`

  - 全局地址固定位置 `/Users/ray/.config/vsi/eslintrc-js.json` // 绝对地址必须 `/` 开头, 否则会被认为是相对 `项目根目录` 的地址. eg: `~/xxx` 会被解析成 `project/~/xxx` 地址.
  - 本地地址固定位置 `eslintrc-js.yml` // 相对地址是项目根目录. 无法解析 `${workspaceRoot}` 路径.
  - eslint 工作原理: 首先需要 eslint 命令行工具 `npm install -g eslint`, eslint 会按照上述指定路径查找 `package.json` 文件, 最后按照 `package.json` 查找相关规则. 所以 `package.json` 必须存在. 需要指定 `--prefix <path>` 安装 eslint rules, `npm install -D --prefix <path> <packages>`

- `.vscode/launch.json`

- `.gitignore`

- `.editorconfig` for github tab_size

- `eslintrc.json` - `npm install eslint rules`
  - 根据 `--eslint-local` 设置，安装在不同的地方. (local | global)

- `package.json`
  - 不管有没有 `--jest`, 或者有没有安装 jest, 都写入 jest setting.

- `src/main.js`

- `src/example.test.js`
  - 根据 `--jest` 选择是否写入 test 文件。

dependencies:

- eslint 命令行工具 - `npm install -g eslint`

- eslint-rules - 根据 `--eslint-local` 使用 `npm install -D <packages...>`, or `npm install --prefix <global_path> -D <packages...>`

- jest 命令行工具 - `npm install -g jest`

<br />

## ts - typescript

flags:
- `--eslint-local` - 将 eslint 安装在本地.
- `--jest` - 安装 jest.

files:
- `.vscode/settings.json` - 同 `js` 设置.

- `.vscode/launch.json`

- `.vscode/tasks.json` - tsc 将 ts 转成 js

- `tsconfig.json` - 比 js 多了 tsconfig 编译设置. 指定了将 ts 编译成 js 的规则.

- `.gitignore`

- `eslintrc.json` - 同 `js` 设置, 只是需要安装的 eslint rules 不同

- `src/main.js`

- `.editorconfig` for github tab_size

- `jest` 设置 - `npm install -D ts-jest @types/jest`; js 不需要安装

- `package.json` - jest setting
  - 根据 `--jest` 决定是否向 package.json 中写入 jest setting. 这里和 `js` 不同，`js` 是不管有没有 `--jest` 都直接写入 jest setting.

- `src/example.test.ts` - jest

dependencies:

- eslint 命令行工具 - `npm install -g eslint`

- eslint-rules - 根据 `--eslint-local` 使用 `npm install -D <packages...>`, or `npm install --prefix <global_path> -D <packages...>`

- jest 命令行工具 - `npm install -g jest`

- ts 可用的 jest 类型 - `npm install -D ts-jest @types/jest`

<br />

## react - ts

flags:
- `--eslint-local` - 将 eslint 安装在本地.

files: - CRA (create react app) 会创建很多文件.

- `.vscode/settings.json` - 不需要 launch.json

- `tsconfig.json` - CRA 自动生成. 可以添加某些设置:
  - `"noImplicitReturns": true` 强制定义出参类型

- `eslint.json` - `npm install -D <eslint rules>`

- `.editorconfig` for github tab_size

dependencies:

- eslint 命令行工具 - `npm install -g eslint`

- eslint-rules - 根据 `--eslint-local` 使用 `npm install -D <packages...>`, or `npm install --prefix <global_path> -D <packages...>`

settings:
- 修改 `package.json` 主要是添加 "rules", 防止 auto compile 的时候因为 unused rule 而报错.

```json
{
  ...
  "eslintConfig": {
    "extends": [...],
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

<br />

## python ? - 测试用 (Beta)

- `.vscode/settings.json`

- `.vscode/launch.json`

- `.gitignore`

<br />
