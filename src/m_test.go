package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// test Get Exec Path Dir
func TestCurrentDirPath(t *testing.T) {
	t.Log(os.Args[0])
	t.Log(filepath.Dir(os.Args[0]))

	// Abs replace relative path to Absolute Path
	curPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(curPath)

	// Abs replace relative path to Absolute Path
	pathDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pathDir)
}

// test Read Dir Not Exist
func TestReadDirNotExist(t *testing.T) {
	_, err := ioutil.ReadDir("/abc")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Error(err)
		return
	} else if errors.Is(err, os.ErrNotExist) {
		t.Log("not exist")
	}
}

// test create file & write file
func TestReadFileNotExist(t *testing.T) {
	f, err := os.OpenFile("/Users/ray/Desktop/newdir/test1.txt", os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		t.Error(err)
		return
	}

	if fi.Size() != 0 {
		t.Log("file not empty!, size: ", fi.Size())
		return
	}

	n, err := f.Write([]byte("abc"))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("write: ", n)
}

// dir file permit at least 755
// if dir exist no error
func TestCreateDirAll(t *testing.T) {
	err := os.MkdirAll("/Users/ray/Desktop/newdir/test1.txt", 0755)
	if err != nil {
		t.Error(err)
		return
	}
}

// if dir exist error
func TestCreateDir(t *testing.T) {
	t.Log("start")
	err := os.Mkdir("/Users/ray/Desktop/newdir", 0755)
	if err != nil {
		t.Error(err)
		return
	}
}
