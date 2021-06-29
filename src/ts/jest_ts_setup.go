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
	"encoding/json"
	"errors"
	"io"

	"os"

	"local/src/util"
)

const testFolder = "test"

// jestFileContent add example of unit test
var jestFileContent = util.FileContent{
	Path:    testFolder + "/example.test.ts",
	Content: exampleTestTS,
}

// TS 中 jest 所需要的依赖
var jestDependencies = []string{"@types/jest", "ts-jest"}

// 查看 package.json devDependencies 是否下载了 @types/jest, ts-jest
// npm i -D @types/jest ts-jest
func dependenciesNeedsToInstall(dependencies []string) (libs []string, err error) {
	// open package.json 文件
	pkgFile, err := os.OpenFile("package.json", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	defer pkgFile.Close()

	// 获取 file info
	pkgInfo, err := pkgFile.Stat()
	if err != nil {
		return nil, err
	}

	// package.json is empty
	if pkgInfo.Size() == 0 {
		// NOTE package.json shouldn't be empty
		return nil, errors.New("package.json shouldn't be empty, please re-initialize the project")
	}

	pkgMap, err := readFileToMap(pkgFile)
	if err != nil {
		return nil, err
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	return _checkDependencies(pkgMap, dependencies)
}

func readFileToMap(packageFile *os.File) (map[string]interface{}, error) {
	byt, err := io.ReadAll(packageFile)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(byt, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// 检查 devDependencies 是否有安装 "ts-jest", "@types/jest"
func _checkDependencies(pkgMap map[string]interface{}, libs []string) ([]string, error) {
	var result []string

	devDependencies, ok := pkgMap["devDependencies"]
	if !ok {
		return libs, nil
	}

	dev, ok := devDependencies.(map[string]interface{})
	if !ok {
		return nil, errors.New("devDependencies assert error: is not an Object")
	}

	// 检查 dependencies 是否存在
	for _, lib := range libs {
		if _, ok := dev[lib]; !ok {
			result = append(result, lib)
		}
	}

	return result, nil
}
