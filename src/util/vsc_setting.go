package util

import (
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

func ReadVscFile() (*os.File, error) {
	vscPath, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	return os.Open(vscPath)
}
