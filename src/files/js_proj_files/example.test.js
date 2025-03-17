// https://jestjs.io/docs/getting-started
// 1. 安装 jest 命令行工具, `npm install --save-dev jest`
// 2. package.json 中加入
//    {
//      "scripts": {
//        "test": "jest"
//      }
//    }
// 3. 测试文件必须以 *.test.js 结尾

// 测试方法:
// npm test                 // 所有文件
// npm test test/*.test.js  // 指定文件
// npm test --coverage {pattern}  // 显示 coverage
// 或者 vscode debug 中选择的 Jest Current File

function add(a, b) {
  return a + b;
}

test('add function test', () => {
  expect(add(1, 2)).toBe(3);
});
