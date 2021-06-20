package util

import (
	"errors"
	"fmt"
	"os"
)

func WriteCfgFiles(folders []string, fileContent map[string][]byte) {
	// create folders
	for _, v := range folders {
		err := createDir(v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// write files
	for fp, fc := range fileContent {
		err := createAndWriteFile(fp, fc)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
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
	defer func() {
		if er := f.Close(); er != nil {
			fmt.Println("file close error:", er)
			return
		}
	}()

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
