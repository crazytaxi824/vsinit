package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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
	fmt.Printf("creating directories: %s ... ", folderPath)
	err := os.Mkdir(folderPath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) { // 判断 dir 是否已经存在
		fmt.Println("failed")
		return fmt.Errorf("create %s Dir error: %w", folderPath, err)
	} else if errors.Is(err, os.ErrExist) {
		// 如果文件夹已经存在
		fmt.Println("skip, already exists")
		return nil
	}

	fmt.Println("done")
	return nil
}

// create and write files.
func createAndWriteFile(fc FileContent) error {
	fmt.Printf("writing file: %s ... ", fc.Path)
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
		fmt.Println("skip, already exists")
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

// npm install libs to devDependencies
// 指定位置安装 eslint 所需依赖
func NpmInstallDependencies(path string, libs ...string) error {
	if len(libs) == 0 {
		return nil
	}

	// TODO 是否需要安装？(y/n)

	results := []string{"i", "-D"}

	// 指定下载到什么地方
	if path != "" {
		results = append(results, "--prefix", path)
	}

	// 执行命令
	results = append(results, libs...)
	cmd := exec.Command("npm", results...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NpmInstallGlobalDependencies(libs ...string) error {
	if len(libs) == 0 {
		return nil
	}

	// TODO 是否需要安装？(y/n)

	results := []string{"i", "-g"}

	// 执行命令
	results = append(results, libs...)
	cmd := exec.Command("npm", results...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Warn(m string) string {
	return fmt.Sprintf("\033[0;37;41m%s\033[0m", m)
}
