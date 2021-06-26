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