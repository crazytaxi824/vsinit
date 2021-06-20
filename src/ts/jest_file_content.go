// 使用 jest 需要安装 jest 命令行工具, npm i jest -g
// 项目中安装 npm i -D @types/jest ts-jest
// 测试文件必须以 *.test.ts 结尾

// 测试方法:
// npm run test /test/example.test.ts
// npm run test-c /test/example.test.ts
// vscode debug 中选择的 Jest Current File

package ts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"local/src/util"
	"os"
	"os/exec"
)

// 项目根目录下生成 test 文件夹
const JestFolder = "test"

// 写入测试命令到 /test/example.test.ts 文件中
var jestExampleTestFile = []byte(`function add(a: number, b: number): number {
  return a + b;
}

test('add function test', () => {
  expect(add(1, 2)).toBe(3);
});
`)

var JestFileContent = util.FileContent{
	Path:    "test/example.test.ts",
	Content: jestExampleTestFile,
}

// package.json 文件中需要用到的设置
var jestPackageConfig = map[string]map[string]string{
	"scripts": {
		"build":  "tsc",
		"test":   "jest",
		"test-c": "jest --coverage",
	},
	"jest": {
		"testEnvironment": "node",
		"preset":          "ts-jest",
	},
}

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
// add "scripts" and "jest" to package.json file
func JestSetup() error {
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
	packageMap, err := readFileToMap(packageFile)
	if err != nil {
		return err
	}

	// 检查 package.json 中是否有 "scripts" 和 "jest"
	err = checkPackageFile(packageMap)
	if err != nil {
		return err
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	err = checkDependencies(packageMap)
	if err != nil {
		return err
	}

	// 清空 package.json 文件写入新内容
	err = truncAndReWrieFile(packageFile, packageMap)
	if err != nil {
		return err
	}

	return nil
}

// npm install ts-jest @types/jest
func npmInstallDependencies(libs ...string) error {
	for _, lib := range libs {
		cmd := exec.Command("npm", "i", "-D", lib)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

// 将文件内容 json 反序列化到 map 中。
func readFileToMap(packageFile *os.File) (map[string]interface{}, error) {
	// 读取 package.json 内容
	packageContent, err := io.ReadAll(packageFile)
	if err != nil {
		return nil, err
	}

	// json 反序列化
	var packageMap map[string]interface{}
	err = json.Unmarshal(packageContent, &packageMap)
	if err != nil {
		return nil, err
	}

	return packageMap, nil
}

// package.json 没有任何内容的情况
func newPackageFile(packageFile *os.File) error {
	// 先将 "scripts" and "jest" 写入 packagefile
	result, err := jsonIndentContent(jestPackageConfig)
	if err != nil {
		return err
	}

	_, err = packageFile.Write(result)
	if err != nil {
		return err
	}

	// npm install ts-jest @types/jest
	err = npmInstallDependencies("ts-jest", "@types/jest")
	if err != nil {
		return err
	}

	return nil
}

// 清空文件，重新写入内容
func truncAndReWrieFile(packageFile *os.File, content map[string]interface{}) error {
	// 将读写符重置到文件的起始位置
	_, err := packageFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 清空文件
	err = packageFile.Truncate(0)
	if err != nil {
		return err
	}

	// json 序列化 & 格式化
	result, err := jsonIndentContent(content)
	if err != nil {
		return err
	}

	// 写入新内容
	_, err = packageFile.Write(result)
	if err != nil {
		return err
	}

	return nil
}

// 检查是否有安装 "ts-jest", "@types/jest"
func checkDependencies(packageMap map[string]interface{}) error {
	devDependencies, ok := packageMap["devDependencies"]
	if !ok {
		err := npmInstallDependencies("ts-jest", "@types/jest")
		if err != nil {
			return err
		}
	}

	devDependenciesMap, ok := devDependencies.(map[string]interface{})
	if !ok {
		return errors.New("devDependencies assert failed")
	}

	if _, ok := devDependenciesMap["ts-jest"]; !ok {
		// download ts-jest
		err := npmInstallDependencies("ts-jest")
		if err != nil {
			return err
		}
	}

	if _, ok := devDependenciesMap["@types/jest"]; !ok {
		// download @types/jest
		err := npmInstallDependencies("@types/jest")
		if err != nil {
			return err
		}
	}
	return nil
}

// 添加修改 package.json 中的 "scripts", "jest" 字段
func checkPackageFile(packageMap map[string]interface{}) error {
	// 修改 "jest" 字段
	err := changePackageFile(packageMap, "jest")
	if err != nil {
		return err
	}

	// 修改 "scripts" 字段
	err = changePackageFile(packageMap, "scripts")
	if err != nil {
		return err
	}

	return nil
}

// 将字段添加到 packageMap
func changePackageFile(packageMap map[string]interface{}, key string) error {
	// 判断 key 是否存在
	if value, ok := packageMap[key]; !ok {
		// key 不存在
		packageMap[key] = jestPackageConfig[key]
	} else {
		// key 存在
		valMap, okk := value.(map[string]interface{})
		if !okk {
			return fmt.Errorf("%s assert failed", key)
		}

		// copy jestPackageConfig to jestMap
		for k, v := range jestPackageConfig[key] {
			valMap[k] = v
		}
	}
	return nil
}

// json marshal & indent file
func jsonIndentContent(v interface{}) ([]byte, error) {
	src, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = json.Indent(&buf, src, "", "  ")
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
