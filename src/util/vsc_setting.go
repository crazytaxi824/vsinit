package util

import (
	"encoding/json"
	"errors"
	"os"
)

type VscSetting struct {
	Golangci string `json:"golangci,omitempty"`
	Eslint   string `json:"eslint,omitempty"`
}

func getConfigDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("$HOME is not exist, please set $HOME env")
	}

	return home + "/.vsc", nil
}

func ReadVscFile() (*VscSetting, *Suggestion, error) {
	vscPath, err := getConfigDir()
	if err != nil {
		return nil, nil, err
	}

	// 检查 "~/.vsc/vsc-config.json" 文件，看是否存在 golangci-lint 配置文件位置。
	vscf, err := os.Open(vscPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, &Suggestion{
			Problem:  "haven't setup golangci-lint yet, please set it:",
			Solution: GolintciCmd,
		}, nil
	}
	defer vscf.Close()

	// json 反序列化
	var cfg VscSetting
	de := json.NewDecoder(vscf)
	err = de.Decode(&cfg)
	if err != nil {
		return nil, nil, err
	}

	return &cfg, nil, nil
}
