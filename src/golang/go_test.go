package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// 测试获取绝对路径
func Test_AbsPath(t *testing.T) {
	fpath := "abc/def"
	t.Log(filepath.Abs(fpath))
}

// 测试获取环境变量
func Test_CheckShell(t *testing.T) {
	t.Log(os.Getenv("SHELL"))
	t.Log(os.Getenv("GOBIN"))
	t.Log(os.Getenv("GOPATH"))
	t.Log(os.Getenv("HOME"))
	if os.Getenv("abcdef") != "" {
		t.Fail()
	}
}

func Test_StringFormat(t *testing.T) {
	err := checkGOPATH()
	if err != nil {
		fmt.Println(err)
	}
}

func Test_replaceLintConfig(t *testing.T) {
	ns, sug, err := replaceLintConfig(golangciConfig, "abc/def.yml")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(sug)
	t.Log(string(ns))
}

func Test_replaceLintConfig2(t *testing.T) {
	ns, sug, err := replaceLintConfig([]byte(""), "abc/def.yml")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(sug)
	t.Log(string(ns))
}

func Test_replaceHolder(t *testing.T) {
	ns, _, err := replaceLintConfig(golangciConfig, "abc/def.yml")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(replaceLintPlaceHolder(ns)))
}

func Test_replaceHolderNil(t *testing.T) {
	t.Log(string(replaceLintPlaceHolder(nil)))
}
