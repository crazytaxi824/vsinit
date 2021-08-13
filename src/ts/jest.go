// 使用 jest 单元测试需要安装 jest 命令行工具, npm i jest -g
// 项目中安装 npm i -D @types/jest ts-jest
// 测试文件必须以 *.test.ts 结尾

// 测试方法:
// npm run test *.test.ts
// npm run test-c *.test.ts
// npm run test-c **/*.test.ts  测试所有 test 文件
// 或者 vscode debug 中选择的 Jest Current File

package ts

import (
	"errors"

	"local/src/util"
)

const testFolder = "test"

// jestFileContent 添加 example of unit test
var jestFileContent = util.FileContent{
	Path:    testFolder + "/example.test.ts",
	Content: exampleTestTS,
}

// TS 中 jest 所需要的依赖
var jestDependencies = []string{"@types/jest", "ts-jest"}

// 写入 Jest 相关文件，test/example.test.ts 文件. 添加 Jest 所需依赖.
func initJest(ctx *util.VSContext) error {
	// 检查 npm 是否安装，把 suggestion 当 error 返回，因为必须要安装依赖
	if sugg := util.CheckCMDInstall("npm"); sugg != nil {
		return errors.New(sugg.String())
	}

	ctx.AddFolders(testFolder)
	ctx.AddFiles(jestFileContent)

	// 在本地添加 Jest 依赖
	return ctx.AddMissingDependencies(jestDependencies, "package.json", "")
}
