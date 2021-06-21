// 使用 jest 单元测试需要安装 jest 命令行工具, npm i jest -g
// 项目中安装 npm i -D @types/jest ts-jest
// 测试文件必须以 *.test.ts 结尾

// 测试方法:
// npm run test xxx.test.ts
// npm run test-c xxx.test.ts
// npm run test-c **/*.test.ts  测试所有 test 文件
// 或者 vscode debug 中选择的 Jest Current File

function add(a: number, b: number): number {
  return a + b;
}

test("add function test", () => {
  expect(add(1, 2)).toBe(3);
});
