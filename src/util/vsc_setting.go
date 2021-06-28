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

	// eslint 文件夹
	EslintDirector = "/eslint" // TODO 移动到 js/ts 文件中
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

func (vs *VscSetting) ReadFromFile(vscDir string) error {
	// read vsc config file
	f, err := os.Open(vscDir + VscConfigFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// ~/.vsc/vsc-config 文件存在, 读取文件
	err = vs.readJSON(f)
	if err != nil {
		return err
	}

	return nil
}

func (vs *VscSetting) readJSON(reader io.Reader) error {
	de := json.NewDecoder(reader)
	return de.Decode(vs)
}

func (vs *VscSetting) JSONIndentFormat() ([]byte, error) {
	return json.MarshalIndent(vs, "", "  ")
}

func GetVscConfigDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("$HOME is not exist, please set $HOME env")
	}

	return home + vscDirectory, nil
}
