package util

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

type FoldersAndFiles struct {
	folders []string
	files   []FileContent
	tsjs    struct {
		dependencies []dependenciesInstall
		suggestions  []*Suggestion
		lintPath     string
	}
}

func InitFoldersAndFiles(folders []string, files []FileContent) *FoldersAndFiles {
	var ff FoldersAndFiles
	ff.folders = folders
	ff.files = files

	return &ff
}

func (ff *FoldersAndFiles) SetLintPath(lintPath string) {
	ff.tsjs.lintPath = lintPath
}

func (ff *FoldersAndFiles) LintPath() string {
	return ff.tsjs.lintPath
}

func (ff *FoldersAndFiles) AddFiles(files ...FileContent) {
	ff.files = append(ff.files, files...)
}

func (ff *FoldersAndFiles) AddFolders(folders ...string) {
	ff.folders = append(ff.folders, folders...)
}

func (ff *FoldersAndFiles) AddSuggestions(sug ...*Suggestion) {
	ff.tsjs.suggestions = append(ff.tsjs.suggestions, sug...)
}

func (ff *FoldersAndFiles) Suggestions() []*Suggestion {
	if len(ff.tsjs.suggestions) > 0 {
		return ff.tsjs.suggestions
	}

	return nil
}

func (ff *FoldersAndFiles) _addDependencies(dependencies ...dependenciesInstall) {
	ff.tsjs.dependencies = append(ff.tsjs.dependencies, dependencies...)
}

// 添加缺失的依赖
func (ff *FoldersAndFiles) AddMissingDependencies(dependencies []string, packageJSONPath, prefix string) error {
	// 检查本地 package.json 文件
	libs, err := checkMissingdependencies(dependencies, packageJSONPath)
	if err != nil {
		return err
	}

	if len(libs) > 0 {
		ff._addDependencies(dependenciesInstall{
			dependencies: libs,
			prefix:       prefix,
		})
	}

	return nil
}

// 安装所有缺失的依赖
func (ff *FoldersAndFiles) InstallMissingDependencies() error {
	if len(ff.tsjs.dependencies) > 0 {
		for _, dep := range ff.tsjs.dependencies {
			err := npmInstallDependencies(dep.prefix, dep.global, dep.dependencies...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 查看 package.json 是否下载了所需要的依赖.
//  - package.json 可以是 local 也可以是 global，需要手动填写文件地址.
func checkMissingdependencies(dependencies []string, packageJSONPath string) (libs []string, err error) {
	// open package.json 文件
	pkgMap, err := _readPkgJSONToMap(packageJSONPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// package.json 不存在的情况，下载所有 dependencies
		return dependencies, nil
	}

	// 查看 devDependencies 是否有下载
	// npm install ts-jest @types/jest
	return _filterDependencies(pkgMap, dependencies)
}

// 读取 package.json 文件, json 反序列化到 map 中.
func _readPkgJSONToMap(packageJSONPath string) (map[string]interface{}, error) {
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
func _filterDependencies(pkgMap map[string]interface{}, libs []string) ([]string, error) {
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

// 写入所需文件
func (ff *FoldersAndFiles) WriteAllFiles() error {
	return WriteFoldersAndFiles(ff.folders, ff.files)
}

// 生成 lint 配置文件，记录 lint 配置文件地址。
func (ff *FoldersAndFiles) AddLintConfigAndLintPath(lintPath string, lincCfgFile []byte) {
	ff.AddFiles(FileContent{
		Path:    lintPath,
		Content: lincCfgFile,
	})

	// eslintrc-ts.json 的文件路径
	ff.tsjs.lintPath = lintPath
}

// 读取 .vscode/settings.json 文件, 获取想要的值
func ReadSettingJSON(v interface{}) error {
	// 读取 .vscode/settings.json
	settingsPath, err := filepath.Abs(SettingsJSONPath)
	if err != nil {
		return err
	}

	sf, err := os.Open(settingsPath)
	if err != nil {
		return err
	}
	defer sf.Close()

	// json 反序列化 settings.json
	jsonc, err := io.ReadAll(sf)
	if err != nil {
		return err
	}

	js, err := JSONCToJSON(jsonc)
	if err != nil {
		return err
	}

	err = json.Unmarshal(js, v)
	if err != nil {
		return err
	}

	return nil
}
