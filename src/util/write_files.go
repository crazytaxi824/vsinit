package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// 需要写入项目的文件
type FileContent struct {
	Filepath string // 文件路径: 包括路径和文件名
	Content  []byte // 文件内容
}

func (fc *FileContent) createDir() error {
	// parse filepath
	dir := filepath.Dir(fc.Filepath)

	// create dir, if dir == "." skip
	if dir != "." {
		err := os.MkdirAll(dir, 0750)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func WriteFiles(fileContents []FileContent) error {
	fmt.Printf("%sWriting files ...%s\n", COLOR_GREEN, COLOR_RESET)

	for _, fc := range fileContents {
		var (
			f   *os.File
			err error
		)

		// mkdir
		err = fc.createDir()
		if err != nil {
			log.Println(err)
			return err
		}

		// write file content
		f, err = os.OpenFile(fc.Filepath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil && !errors.Is(err, os.ErrExist) {
			log.Println(err)
			return err
		} else if errors.Is(err, os.ErrExist) {
			// if the file is already exits then skip.
			fmt.Printf("  - %s ... skip, file exists.\n", fc.Filepath)
			continue
		}

		_, err = f.Write(fc.Content)
		if err != nil {
			log.Println(err)
			return err
		}

		fmt.Printf("%s - %s ... done%s\n", COLOR_GREEN, fc.Filepath, COLOR_RESET)
	}

	return nil
}
