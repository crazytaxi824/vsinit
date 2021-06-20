// 使用 jest 需要安装 jest 命令行工具, npm i jest -g
// 项目中安装 npm i -D @types/jest ts-jest
// 测试文件必须以 *.test.ts 结尾

// 测试方法:
// npm run test /test/example.test.ts
// npm run test-c /test/example.test.ts
// vscode debug 中选择的 Jest Current File

package ts

// 项目根目录下生成 test 文件夹

// 读取 package.json 文件

// 查看 package.json devDependencies, dependencies 是否下载了 @types/jest, ts-jest
// npm i -D @types/jest ts-jest

// 查看 package.json 是否写了 "jest", "scripts" 字段
//   "scripts": {
//     "build": "tsc",
//     "test": "jest",
//     "test-c": "jest --coverage"
//   },
//   "jest": {
//     "testEnvironment": "node",
//     "preset": "ts-jest"
//   },

// 写入测试命令到 /test/example.test.ts 文件中
