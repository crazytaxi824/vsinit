package main

import (
	"flag"
	"fmt"
	"os/user"
	"strings"
	"testing"

	"local/src/resource"
)

func Test_Flags(t *testing.T) {
	foo := flag.String("foo", "abc", "foo usage")

	f := flag.Lookup("foo")

	flag.Var(f.Value, "f", fmt.Sprintf("alias to %s", f.Name))

	flag.Parse()

	t.Log(*foo)
}

func Test_Color(t *testing.T) {
	const foo = `abc%s123%sfoo`

	t.Logf(fmt.Sprintf(foo, "\x1b[33m", "\x1b[0m"))
}

func Test_FS(t *testing.T) {
	fs := resource.ReactCommonFuncs

	fsd, err := fs.ReadDir("react_proj_files/react_common_fn")
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range fsd {
		if strings.Contains(v.Name(), ".tsx") {
			b, err := fs.ReadFile("react_proj_files/react_common_fn/" + v.Name())
			if err != nil {
				t.Error(err)
				return
			}

			t.Log(string(b))
		}
	}
}

func Test_FilePath(t *testing.T) {
	// src := "~/.config/coc"

	u, err := user.Current()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(u.HomeDir)
}
