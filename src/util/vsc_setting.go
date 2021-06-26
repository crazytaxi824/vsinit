package util

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

const (
	// vsc 文件夹，~/.vsc
	vscDirectory = "/.vsc"

	// vsc 配置文件, ~/.vsc/vsc-config.json
	VscConfigFilePath = "/vsc-config.json"

	// golangci 文件夹
	GolangciDirector = "/golangci"

	// eslint 文件夹
	EslintDirector = "/eslint"
)

// lint 类型
type Lint byte

const (
	Golangci Lint = 1
	Eslint   Lint = 2
)

// config 文件设置
type VscSetting struct {
	Golangci string `json:"golangci,omitempty"`
	Eslint   string `json:"eslint,omitempty"`
}

func (vs *VscSetting) readJSON(reader io.Reader) error {
	de := json.NewDecoder(reader)
	return de.Decode(vs)
}

func (vs *VscSetting) setLintConfig(lint Lint, cfgPath string) {
	switch lint {
	case Golangci:
		vs.Golangci = cfgPath
	case Eslint:
		vs.Eslint = cfgPath
	}
}

func (vs *VscSetting) writeToFile(file *os.File) error {
	r, err := json.Marshal(vs)
	if err != nil {
		return err
	}

	// 重置 I/O offset
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	// 清空文件
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = file.Write(r)
	if err != nil {
		return err
	}

	return nil
}

func GetVscConfigDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("$HOME is not exist, please set $HOME env")
	}

	return home + vscDirectory, nil
}

func ReadVscConfig() (*VscSetting, error) {
	vscDir, err := GetVscConfigDir()
	if err != nil {
		return nil, err
	}

	vscConfigFile, err := os.Open(vscDir + VscConfigFilePath)
	if err != nil {
		return nil, err
	}
	defer vscConfigFile.Close()

	var vscSetting VscSetting
	err = vscSetting.readJSON(vscConfigFile)
	if err != nil {
		return nil, err
	}

	return &vscSetting, nil
}

func SetVscSetting(lint Lint, cfgPath string) error {
	vscDir, err := GetVscConfigDir()
	if err != nil {
		return err
	}

	vscf, err := os.OpenFile(vscDir+VscConfigFilePath, os.O_RDWR, 0600)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// 文件不存在, 写入新文件
		er := writeNewSettingFile(vscDir, cfgPath, lint)
		if er != nil {
			return er
		}
		return nil
	}
	defer vscf.Close()

	// json 反序列化
	var vscSetting VscSetting
	err = vscSetting.readJSON(vscf)
	if err != nil {
		return err
	}

	// 修改设置然后写入
	vscSetting.setLintConfig(lint, cfgPath)

	// 写入文件
	return vscSetting.writeToFile(vscf)
}

func writeNewSettingFile(dirPath, cfgPath string, lint Lint) error {
	// 创建文件夹
	err := os.Mkdir(dirPath, 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	// 写入文件
	vscf, err := os.OpenFile(dirPath+VscConfigFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer vscf.Close()

	var vscSetting VscSetting
	vscSetting.setLintConfig(lint, cfgPath)

	// 写入文件
	return vscSetting.writeToFile(vscf)
}
