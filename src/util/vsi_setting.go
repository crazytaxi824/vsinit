package util

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	// vsi 文件夹，~/.vsi
	vsiDirectory = "/.vsi"

	// vsi 配置文件, ~/.vsi/vsi-config.json
	VsiConfigFilePath = "/vsi-config.json"
)

// config 文件设置
type VsiConfigJSON struct {
	Golangci string `json:"golangci,omitempty"`
	Eslint   struct {
		TS string `json:"typescript,omitempty"`
		JS string `json:"javascript,omitempty"`
	} `json:"eslint,omitempty"`
}

// 全局 vsi 配置文件地址 ~/.vsi
func GetVsiConfigDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("$HOME is not exist, please set $HOME env")
	}

	return home + vsiDirectory, nil
}

// 从指定 dir 中读取 vsi-config.json 文件
func (vs *VsiConfigJSON) ReadFromDir(vsiDir string) error {
	// read vsi config file
	f, err := os.Open(vsiDir + VsiConfigFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// ~/.vsi/vsi-config 文件存在, 读取文件
	de := json.NewDecoder(f)
	err = de.Decode(vs)
	if err != nil {
		return err
	}

	return nil
}

// json 序列化格式化
func (vs *VsiConfigJSON) JSONIndentFormat() ([]byte, error) {
	return json.MarshalIndent(vs, "", "  ")
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
