package util

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type FileContent struct {
	Path      string
	Content   []byte
	Overwrite bool
}

// create folders and write project files.
func WriteFoldersAndFiles(folders []string, fileContents []FileContent) error {
	// create folders
	for _, v := range folders {
		err := createDir(v)
		if err != nil {
			return err
		}
	}

	// write files
	for _, fc := range fileContents {
		err := createAndWriteFile(fc)
		if err != nil {
			return err
		}
	}
	return nil
}

func createDir(folderPath string) error {
	err := os.Mkdir(folderPath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) { // 判断 dir 是否已经存在
		return fmt.Errorf("create %s Dir error: %w", folderPath, err)
	} else if errors.Is(err, os.ErrExist) {
		// 如果文件夹已经存在
		return nil
	}

	return nil
}

// create and write files.
func createAndWriteFile(fc FileContent) error {
	fmt.Printf("    %s ... ", fc.Path)
	f, err := os.OpenFile(fc.Path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("create %s Files error: %w", fc.Path, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("get %s File status error: %w", fc.Path, err)
	}

	// file is not empty, DO NOT TOUCH. Unless Overwrite
	if fi.Size() != 0 && !fc.Overwrite {
		fmt.Println("skip")
		return nil
	}

	if fc.Overwrite { // 如果重写文件需要 truncate
		if _, er := f.Seek(0, io.SeekStart); er != nil {
			return er
		}

		if er := f.Truncate(0); er != nil {
			return er
		}
	}

	// write file content
	_, err = f.Write(fc.Content)
	if err != nil {
		fmt.Println("failed")
		return fmt.Errorf("write file %s error: %w", fc.Path, err)
	}

	fmt.Println("done")
	return nil
}
