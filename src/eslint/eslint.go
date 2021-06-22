package eslint

import (
	"errors"
	"fmt"
	"io"
	"local/src/util"
	"net/http"
	"os"
)

// github eslint config file addr
const eslintCfgAddr = "https://raw.githubusercontent.com/crazytaxi824/lints/main/eslintrc-ts.json"

var dependencies = []string{
	"eslint-plugin-import",
	"eslint-plugin-jsx-a11y",
	"eslint-plugin-react",
	"eslint-plugin-react-hooks",
	"@typescript-eslint/parser",        // parser
	"@typescript-eslint/eslint-plugin", // main plugin
	"eslint-plugin-jest",               // jest unit test
	"eslint-plugin-promise",            // promise 用法
	"eslint-config-airbnb-typescript",  // ts 用
	"eslint-config-airbnb-base",        // js 专用 lint

	// 解决 prettier 格式化和 eslint 之间的冲突
	"eslint-config-prettier",
}

// TODO 如何判断 eslint 已经成功安装过了？
// NOTE eslintrc.json 必须和 package.json 在同一个文件夹下
func InstallEslintGlobally(fpath string) (eslintCfgFile string, err error) {
	// check eslint command
	err = util.CheckCMDInstall("eslint")
	if err != nil {
		return "", err
	}

	cfgPath := fpath + "/eslint"

	// directory must exist.
	err = os.Mkdir(fpath+"/eslint", 0750)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	} else if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("%s is not exist", fpath)
	}

	// 指定位置安装 eslint 所需依赖
	// NOTE npm install --prefix ./install/here <package>
	err = util.NpmInstallDependencies(cfgPath, dependencies...)
	if err != nil {
		return "", err
	}

	// github 获取配置文件
	cfgbyte, err := githubConfigFile()
	if err != nil {
		return "", nil
	}

	// 写入文件
	cfgPath += "/eslintrc-ts.json"
	err = writeFile(cfgPath, cfgbyte)
	if err != nil {
		return "", nil
	}

	// 返回 eslint config file path
	return cfgPath, nil
}

func githubConfigFile() ([]byte, error) {
	resp, err := http.Get(eslintCfgAddr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func writeFile(fpath string, content []byte) error {
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}
