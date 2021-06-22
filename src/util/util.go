package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type FileContent struct {
	Path    string
	Content []byte
}

// create folders and write project files.
func WriteCfgFiles(folders []string, fileContents []FileContent) error {
	// create folders
	for _, v := range folders {
		err := createDir(v)
		if err != nil {
			return err
		}
	}

	// write files
	for _, fc := range fileContents {
		err := createAndWriteFile(fc.Path, fc.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDir(folderPath string) error {
	fmt.Printf("creating directories: %s ... ", folderPath)
	err := os.Mkdir(folderPath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) { // 判断 dir 是否已经存在
		fmt.Println("failed")
		return fmt.Errorf("create %s Dir error: %w", folderPath, err)
	}

	fmt.Println("done")
	return nil
}

// create and write files.
func createAndWriteFile(fpath string, content []byte) error {
	fmt.Printf("writing file: %s ... ", fpath)
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("create %s Files error: %w", fpath, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("get %s File status error: %w", fpath, err)
	}

	// file is not empty, DO NOT TOUCH.
	if fi.Size() != 0 {
		fmt.Println("skip, file already exists.")
		return nil
	}

	// write file content
	_, err = f.Write(content)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("write file %s error: %w", fpath, err)
	}

	fmt.Println("done")
	return nil
}

// unescape \uxxxx in json string
func UnescapeStringInJSON(src string) (string, error) {
	// FIXME jsonvalue 的问题，等待更新
	// 先处理 \/ 问题
	tmp := strings.Replace(src, `\/`, "/", -1)

	// NOTE 注意 repalce 的时候只能用 `` 符号，否则 \\ 在一起是转义的. 需要用 4 个 \\\\u
	return strconv.Unquote(strings.Replace(strconv.Quote(tmp), `\\u`, `\u`, -1))
}

// npm install libs to devDependencies
func NpmInstallDependencies(libs ...string) error {
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
