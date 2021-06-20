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
	"io"
	"local/src/util"
	"os"
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
func ChangePackageFile() error {
	// 读取 package.json 文件
	packageFile, err := os.OpenFile("package.json", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer packageFile.Close()

	// 读取 file info
	packageInfo, err := packageFile.Stat()
	if err != nil {
		return err
	}

	// package.json is empty
	if packageInfo.Size() == 0 {
		// 先将 "scripts" and "jest" 写入 packagefile
		result, er := jsonIndentContent(jestPackageConfig)
		if er != nil {
			return er
		}

		_, er = packageFile.Write(result)
		if er != nil {
			return er
		}

		// npm install ts-jest @types/jest
		npmInstalldependencies("ts-jest", "@types/jest")

		return nil
	}

	// package.json is not empty
	// 读取 package.json 内容
	packageContent, err := io.ReadAll(packageFile)
	if err != nil {
		return err
	}

	// json 反序列化
	var packageMap map[string]interface{}
	err = json.Unmarshal(packageContent, &packageMap)
	if err != nil {
		return err
	}

	// 判断 "scripts" 是否存在
	if scripts, ok := packageMap["scripts"]; !ok {
		// "scripts" 不存在
		packageMap["scripts"] = jestPackageConfig["scripts"]
	} else {
		// "scripts" 存在
		scriptsMap, okk := scripts.(map[string]interface{})
		if !okk {
			return errors.New("script assert failed")
		}

		// copy jestPackageConfig to scriptsMap
		for k, v := range jestPackageConfig["scripts"] {
			scriptsMap[k] = v
		}
	}

	// 判断 "jest" 是否存在
	if jest, ok := packageMap["jest"]; !ok {
		// "jest" 不存在
		packageMap["jest"] = jestPackageConfig["jest"]
	} else {
		// "jest" 存在
		jestMap, okk := jest.(map[string]interface{})
		if !okk {
			return errors.New("jest assert failed")
		}

		// copy jestPackageConfig to jestMap
		for k, v := range jestPackageConfig["jest"] {
			jestMap[k] = v
		}
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	devDependencies, ok := packageMap["devDependencies"]
	if !ok {
		npmInstalldependencies("ts-jest", "@types/jest")
	}

	devDependenciesMap, ok := devDependencies.(map[string]interface{})
	if !ok {
		return errors.New("devDependencies assert failed")
	}

	if _, ok := devDependenciesMap["ts-jest"]; !ok {
		// download ts-jest
		npmInstalldependencies("ts-jest")
	}

	if _, ok := devDependenciesMap["@types/jest"]; !ok {
		// download @types/jest
		npmInstalldependencies("@types/jest")
	}

	// 将读写符重置到文件的起始位置
	_, err = packageFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 清空文件
	err = packageFile.Truncate(0)
	if err != nil {
		return err
	}

	// 写入新内容
	result, err := jsonIndentContent(packageMap)
	if err != nil {
		return err
	}

	_, err = packageFile.Write(result)
	if err != nil {
		return err
	}

	return nil
}

// TODO npm install ts-jest @types/jest
func npmInstalldependencies(lib ...string) {}

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
