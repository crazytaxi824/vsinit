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

// 写入 test 相关文件，test/example.test.ts 文件
func (ff *foldersAndFiles) initJest() error {
	// 检查 npm 是否安装，把 suggestion 当 error 返回，因为必须要安装依赖
	if sugg := util.CheckCMDInstall("npm"); sugg != nil {
		return errors.New(sugg.String())
	}

	ff._addFolders(testFolder)
	ff._addFiles(jestFileContent)

	// 添加 Jest Dependencies
	return ff.addMissingJestDependencies()
}

func (ff *foldersAndFiles) addMissingJestDependencies() error {
	// 检查本地 package.json 文件
	libs, err := checkMissingdependencies(jestDependencies, "package.json")
	if err != nil {
		return err
	}

	if len(libs) > 0 {
		ff._addDependencies(util.DependenciesInstall{
			Dependencies: libs,
			Prefix:       "",
			Global:       false,
		})
	}

	return nil
}

// 查看 package.json devDependencies 是否下载了 @types/jest, ts-jest
// npm i -D @types/jest ts-jest
func checkMissingdependencies(dependencies []string, pkgFilePath string) (libs []string, err error) {
	// open package.json 文件
	pkgFile, err := os.Open(pkgFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// package.json 不存在的情况，下载所有 dependencies
		return dependencies, nil
	}
	defer pkgFile.Close()

	pkgMap, err := _readPkgJSONToMap(pkgFile)
	if err != nil {
		return nil, err
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	return _filterDependencies(pkgMap, dependencies)
}

func _readPkgJSONToMap(packageFile *os.File) (map[string]interface{}, error) {
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

// 筛选 devDependencies 是否有安装 "ts-jest", "@types/jest"
func _filterDependencies(pkgMap map[string]interface{}, libs []string) ([]string, error) {
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
