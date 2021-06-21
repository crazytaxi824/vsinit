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
	_ "embed" // for go:embed file use
	"errors"
	"os"
	"os/exec"

	"local/src/jest"
	"local/src/util"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

var (
	//go:embed cfgfiles/example.test.ts
	exampleTestFile []byte

	//go:embed cfgfiles/packagecfg.json
	packageCfgJSON []byte
)

var TSJestFileContent = util.FileContent{
	Path:    "test/example.test.ts",
	Content: exampleTestFile,
}

// 查看 package.json devDependencies, dependencies 是否下载了 @types/jest, ts-jest
// npm i -D @types/jest ts-jest
func SetupTS() error {
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
		// 删除后重新下载依赖
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

// package.json 没有任何内容的情况下直接写文件
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

// 添加修改 package.json 中的 "scripts", "jest" 字段
func setPackageFile(packageFile *os.File, packageRootV *jsonvalue.V) error {
	// 反序列化 package.json 配置文件内容
	packageConfig, err := jsonvalue.Unmarshal(packageCfgJSON)
	if err != nil {
		return err
	}

	// 修改 "jest" 字段
	err = jest.CheckPackageFile(packageRootV, packageConfig, "jest")
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
