// 使用 jest 单元测试需要安装 jest 命令行工具, npm i jest -g
// 测试文件必须以 *.test.js 结尾

// 测试方法:
// npm run test *.test.js
// npm run test-c *.test.js
// npm run test-c **/*.test.js  测试所有 test 文件
// 或者 vscode debug 中选择的 Jest Current File

package js

import (
	_ "embed" // for go:embed file use

	"local/src/jest"
	"local/src/util"
	"os"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

var (
	//go:embed cfgfiles/example.test.js
	exampleTestFile []byte

	//go:embed cfgfiles/packagecfg.json
	packageCfgJSON []byte
)

var JestFileContent = util.FileContent{
	Path:    "test/example.test.js",
	Content: exampleTestFile,
}

func SetupJS() error {
	// open package.json 文件
	packageFile, err := os.OpenFile("package.json", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer packageFile.Close()

	// 获取 file info
	packageInfo, err := packageFile.Stat()
	if err != nil {
		return err
	}

	// package.json is empty
	if packageInfo.Size() == 0 {
		err = newPackageFile(packageFile)
		if err != nil {
			return err
		}
		return nil
	}

	// package.json is not empty
	// 反序列化读取 package.json 配置文件
	packageRootV, err := jest.ReadFileToJsonvalue(packageFile)
	if err != nil {
		return err
	}

	// 检查 package.json 中是否有 "scripts"
	err = setPackageFile(packageFile, packageRootV)
	if err != nil {
		return err
	}

	return nil
}

func newPackageFile(packageFile *os.File) error {
	_, err := packageFile.Write(packageCfgJSON)
	if err != nil {
		return err
	}

	return nil
}

// 添加修改 package.json 中的 "scripts" 字段
func setPackageFile(packageFile *os.File, packageRootV *jsonvalue.V) error {
	// 反序列化 package.json 配置文件内容
	packageConfig, err := jsonvalue.Unmarshal(packageCfgJSON)
	if err != nil {
		return err
	}

	// 修改 "scripts" 字段
	err = jest.CheckPackageFile(packageRootV, packageConfig, "scripts")
	if err != nil {
		return err
	}

	// 清空 package.json 文件写入新内容
	err = jest.WrieFile(packageFile, packageRootV)
	if err != nil {
		return err
	}

	return nil
}
