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
	_ "embed" // for go:embed file use
	"encoding/json"
	"errors"
	"io"
	"local/src/util"
	"os"
	"os/exec"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// 项目根目录下生成 test 文件夹
const JestFolder = "test"

var (
	//go:embed jesttestfiles/example.test.ts
	jestExampleTestFile []byte

	//go:embed jesttestfiles/packagecfg.json
	packageCfgJSON []byte
)

var JestFileContent = util.FileContent{
	Path:    "test/example.test.ts",
	Content: jestExampleTestFile,
}

// 查看 package.json devDependencies, dependencies 是否下载了 @types/jest, ts-jest
// npm i -D @types/jest ts-jest
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
	packageRootV, err := readFileToJsonvalue(packageFile)
	if err != nil {
		return err
	}

	// 检查 package.json 中是否有 "scripts" 和 "jest"
	err = setPackageFile(packageFile, packageRootV)
	if err != nil {
		return err
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	err = checkDependencies(packageRootV)
	if err != nil {
		return err
	}

	return nil
}

// 检查 devDependencies and dependencies 是否有安装 "ts-jest", "@types/jest"
func checkDependencies(packageRootV *jsonvalue.V) error {
	// 检查 dependencies 是否存在
	depV, err := packageRootV.Get("devDependencies")
	if err != nil && !errors.Is(err, jsonvalue.ErrNotFound) {
		return err
	} else if errors.Is(err, jsonvalue.ErrNotFound) {
		// devDependencies 不存在, npm install ts-jest @types/jest
		er := npmInstallDependencies("ts-jest", "@types/jest")
		if er != nil {
			return er
		}
		return nil
	} else if !depV.IsObject() {
		// devDependencies 存在, 但不是 object
		er := packageRootV.Delete("devDependencies")
		if er != nil {
			return er
		}

		// npm install ts-jest @types/jest
		er = npmInstallDependencies("ts-jest", "@types/jest")
		if er != nil {
			return er
		}
		return nil
	}

	// devDependencies 存在, 而且是 object
	// 检查 dependencies 中的依赖是否存在 "ts-jest", "@types/jest"
	tsjest, err := checkLib(depV, "ts-jest")
	if err != nil {
		return err
	}

	typeJest, err := checkLib(depV, "@types/jest")
	if err != nil {
		return err
	}

	var libs []string
	if !tsjest {
		libs = append(libs, "ts-jest")
	}

	if !typeJest {
		libs = append(libs, "@types/jest")
	}

	// npm install ts-jest @types/jest
	err = npmInstallDependencies(libs...)
	if err != nil {
		return err
	}

	return nil
}

// 检查是否有安装 "ts-jest", "@types/jest"
func checkLib(depV *jsonvalue.V, lib string) (bool, error) {
	_, err := depV.Get(lib)
	if err != nil && !errors.Is(err, jsonvalue.ErrNotFound) {
		return false, err
	} else if errors.Is(err, jsonvalue.ErrNotFound) {
		return false, nil
	}

	return true, nil
}

// 将文件内容 json 反序列化到 map 中。
func readFileToJsonvalue(packageFile *os.File) (*jsonvalue.V, error) {
	// 读取 package.json 内容
	packageContent, err := io.ReadAll(packageFile)
	if err != nil {
		return nil, err
	}

	// json 反序列化
	return jsonvalue.Unmarshal(packageContent)
}

// package.json 没有任何内容的情况
func newPackageFile(packageFile *os.File) error {
	_, err := packageFile.Write(packageCfgJSON)
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
func truncAndWrieFile(packageFile *os.File, root *jsonvalue.V) error {
	// json 写入字段顺序
	orders := jsonvalue.Opt{
		MarshalKeySequence: []string{"name", "version", "description",
			"main", "scripts", "jest", "keywords", "author", "license", "key"},
	}

	// json 序列化
	marshbytes, err := root.Marshal(orders)
	if err != nil {
		return err
	}

	// 格式化
	var buf bytes.Buffer
	err = json.Indent(&buf, marshbytes, "", "  ")
	if err != nil {
		return err
	}

	// unescape json string
	result, err := util.UnescapeStringInJSON(buf.String())
	if err != nil {
		return err
	}

	return writeFile(packageFile, result)
}

func writeFile(packageFile *os.File, content string) error {
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

	// 写入新内容
	_, err = packageFile.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// 添加修改 package.json 中的 "scripts", "jest" 字段
func setPackageFile(packageFile *os.File, packageRootV *jsonvalue.V) error {
	// 反序列化 package.json 配置文件内容
	packageConfig, err := jsonvalue.Unmarshal(packageCfgJSON)
	if err != nil {
		return err
	}

	// 修改 "jest" 字段
	err = checkPackageFile(packageRootV, packageConfig, "jest")
	if err != nil {
		return err
	}

	// 修改 "scripts" 字段
	err = checkPackageFile(packageRootV, packageConfig, "scripts")
	if err != nil {
		return err
	}

	// 清空 package.json 文件写入新内容
	err = truncAndWrieFile(packageFile, packageRootV)
	if err != nil {
		return err
	}

	return nil
}

// 将字段添加到 packageMap
func checkPackageFile(packageRootV, packageConfig *jsonvalue.V, key string) error {
	// 判断 key 是否存在
	value, err := packageRootV.Get(key)
	if err != nil && !errors.Is(err, jsonvalue.ErrNotFound) {
		return err
	} else if !errors.Is(err, jsonvalue.ErrNotFound) && !value.IsObject() {
		// 如果 jest | scripts 存在，但是不是 object 的情况
		er := packageRootV.Delete(key)
		if er != nil {
			return er
		}
	}

	cfgV, err := packageConfig.Get(key)
	if err != nil {
		return err
	}

	// 插入数据
	cfgV.RangeObjects(func(k string, v *jsonvalue.V) bool {
		_, er := packageRootV.Set(v).At(key, k)
		if er != nil {
			err = er
			return false
		}
		return true
	})
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
