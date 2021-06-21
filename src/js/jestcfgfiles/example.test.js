// 使用 jest 单元测试需要安装 jest 命令行工具, npm i jest -g
// 测试文件必须以 *.test.js 结尾

// 测试方法:
// npm run test  当前文件
// npm run test-c
// npm run test-c test/example.test.ts  指定文件
// npm run test-c **/*.test.ts  测试所有 test 文件
// 或者 vscode debug 中选择的 Jest Current File

function add(a, b) {
  return a + b;
}

test("add function test", () => {
  expect(add(1, 2)).toBe(3);
});
