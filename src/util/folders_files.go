package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// 需要写入项目的文件
type FileContent struct {
	Path      string // 文件地址
	Content   []byte // 文件内容
	Overwrite bool   // 是否需要覆盖文件内容
}

// 在不同情况下添加不同的文件夹和文件，以便于最后统一写文件。
type VSContext struct {
	folders     []string      // 需要创建的文件夹
	files       []FileContent // 需要写入项目的文件
	suggestions []*Suggestion // 需要手动设置的建议
	lintPath    string        // lint 配置文件的地址，golangci-lint, eslint
	tsjs        struct {
		dependencies []dependenciesInstall // 需要安装的 npm 依赖
	}
}

// 初始化 FoldersAndFiles 对象
func InitFoldersAndFiles(folders []string, files []FileContent) *VSContext {
	var ctx VSContext
	ctx.folders = folders
	ctx.files = files

	return &ctx
}

func (ctx *VSContext) SetLintPath(lintPath string) {
	ctx.lintPath = lintPath
}

func (ctx *VSContext) LintPath() string {
	return ctx.lintPath
}

func (ctx *VSContext) AddFiles(files ...FileContent) {
	ctx.files = append(ctx.files, files...)
}

func (ctx *VSContext) AddFolders(folders ...string) {
	ctx.folders = append(ctx.folders, folders...)
}

func (ctx *VSContext) AddSuggestions(sug ...*Suggestion) {
	ctx.suggestions = append(ctx.suggestions, sug...)
}

func (ctx *VSContext) Suggestions() []*Suggestion {
	if len(ctx.suggestions) > 0 {
		return ctx.suggestions
	}

	return nil
}

func (ctx *VSContext) _addDependencies(dependencies ...dependenciesInstall) {
	ctx.tsjs.dependencies = append(ctx.tsjs.dependencies, dependencies...)
}

// 添加缺失的依赖
func (ctx *VSContext) AddMissingDependencies(dependencies []string, packageJSONPath, prefix string) error {
	// 检查本地 package.json 文件
	libs, err := checkMissingdependencies(dependencies, packageJSONPath)
	if err != nil {
		return err
	}

	if len(libs) == 0 {
		return nil
	}

	// NOTE 判断 global & prefix 是否相同，如果相同直接 append 到里面
	for i, v := range ctx.tsjs.dependencies {
		if v.prefix == prefix {
			ctx.tsjs.dependencies[i].dependencies = append(ctx.tsjs.dependencies[i].dependencies, libs...)
			return nil
		}
	}

	// 如果没有相同的 prefix & global 则整个 append.
	ctx._addDependencies(dependenciesInstall{
		dependencies: libs,
		prefix:       prefix,
	})

	return nil
}

// 查看 package.json 是否下载了所需要的依赖.
//  - package.json 可以是 local 也可以是 global，需要手动填写文件地址.
func checkMissingdependencies(dependencies []string, packageJSONPath string) (libs []string, err error) {
	// open package.json 文件
	pkgMap, err := readPkgJSONToMap(packageJSONPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// package.json 不存在的情况，下载所有 dependencies
		return dependencies, nil
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	return filterDependencies(pkgMap, dependencies)
}

// 读取 package.json 文件, json 反序列化到 map 中.
func readPkgJSONToMap(packageJSONPath string) (map[string]interface{}, error) {
	pkgFile, err := os.Open(packageJSONPath)
	if err != nil {
		return nil, err
	}
	defer pkgFile.Close()

	byt, err := io.ReadAll(pkgFile)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(byt, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// 筛选 "devDependencies" 中没有下载的依赖.
func filterDependencies(pkgMap map[string]interface{}, libs []string) ([]string, error) {
	var result []string

	devDependencies, ok := pkgMap["devDependencies"]
	if !ok {
		return libs, nil
	}

	dev, ok := devDependencies.(map[string]interface{})
	if !ok {
		return nil, errors.New("devDependencies assert error: is not an Object")
	}

	// 检查依赖是否存在
	for _, lib := range libs {
		if _, ok := dev[lib]; !ok {
			result = append(result, lib)
		}
	}

	return result, nil
}

// 安装所有缺失的依赖 // TODO 是否提示需要安装？(y/n)
func (ctx *VSContext) InstallMissingDependencies() error {
	if len(ctx.tsjs.dependencies) > 0 {
		for _, dep := range ctx.tsjs.dependencies {
			if dep.prefix == "" {
				fmt.Printf("npm installing following dependencies at Project Root:\n")
			} else {
				fmt.Printf("npm installing following dependencies at %s:\n", dep.prefix)
			}
			fmt.Println("  " + strings.Join(dep.dependencies, "\n  "))

			err := npmInstallDependencies(dep)
			if err != nil {
				return err
			}
		}
		fmt.Println()
	}

	return nil
}

// 生成 lint 配置文件，记录 lint 配置文件地址。
func (ctx *VSContext) AddLintConfigAndLintPath(lintPath string, lincCfgFile []byte) {
	ctx.AddFiles(FileContent{
		Path:    lintPath,
		Content: lincCfgFile,
	})

	// eslintrc-ts.json 的文件路径
	ctx.lintPath = lintPath
}

// 写入项目所需文件
func (ctx *VSContext) WriteAllFiles() error {
	fmt.Println("writing file: ")

	err := writeFoldersAndFiles(ctx.folders, ctx.files)
	if err != nil {
		return err
	}

	fmt.Println()
	return nil
}

// create folders and write project files.
func writeFoldersAndFiles(folders []string, fileContents []FileContent) error {
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
	fmt.Printf("  %s ... ", fc.Path)
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
