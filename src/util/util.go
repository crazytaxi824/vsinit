package util

import (
	"errors"
	"fmt"
	"os"
)

func WriteCfgFiles(fileContent map[string]string) error {
	// create .vscode & src Dir
	fmt.Printf("creating .vscode & src directories ... ")
	err := createVsCodeDirs()
	if err != nil {
		fmt.Println("fail")
		return err
	}
	fmt.Println("done")

	for fp, fc := range fileContent {
		err = createAndWriteFiles(fp, fc)
		if err != nil {
			return err
		}
	}
	return nil
}

// create .vscode & src dir,
func createVsCodeDirs() error {
	err := os.Mkdir(".vscode", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("create .vscode Dir error: %w", err)
	}

	err = os.Mkdir("src", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("create src Dir error: %w", err)
	}

	return nil
}

// create and write files.
func createAndWriteFiles(fpath, content string) error {
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("create %s Files error: %w", fpath, err)
	}
	defer func() {
		if er := f.Close(); er != nil {
			fmt.Println(er)
			return
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("get %s File status error: %w", fpath, err)
	}

	// file is not empty, DO NOT TOUCH.
	if fi.Size() != 0 {
		return nil
	}

	fmt.Printf("writing file: %s ... ", fpath)
	// write file content
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println("fail")
		return fmt.Errorf("write file %s error: %w", fpath, err)
	}
	fmt.Println("done")

	return nil
}
