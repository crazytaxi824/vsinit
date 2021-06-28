// 使用 jest 单元测试需要安装 jest 命令行工具, npm i jest -g
// 项目中安装 npm i -D @types/jest ts-jest
// 测试文件必须以 *.test.ts 结尾

// 测试方法:
// jest                 // 所有文件
// jest test/*.test.js  // 指定文件
// jest --coverage xxx  // 显示 coverage
// 或者 vscode debug 中选择的 Jest Current File

function add(a: number, b: number): number {
  return a + b;
}

test('add function test', () => {
  expect(add(1, 2)).toBe(3);
});
