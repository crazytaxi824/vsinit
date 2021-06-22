package eslint

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
)

func Test_MkNoneExistDir(t *testing.T) {
	err := os.Mkdir("abc/def/gh", 0750)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Error(err)
		return
	} else if errors.Is(err, os.ErrNotExist) {
		t.Log("ok")
	}
}

func Test_GetConfigFile(t *testing.T) {
	resp, err := http.Get("https://raw.githubusercontent.com/crazytaxi824/lints/main/eslintrc-ts.json")
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(b))
}
