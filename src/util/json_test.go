package util

import (
	"strings"
	"testing"
)

func Test_FindLastChar(t *testing.T) {
	src := `"s": "abc\" def\" \\\t \n /abc //notComment { } ~/omg/gg.json", // comment`

	lastIndex, hasComments, err := analyseJSONCstatement([]byte(src), 0)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(src[:lastIndex+1])
	t.Log(hasComments)
}

const totalsrc = `{
  // this is comment
  "s": "abc\" def\" \\\t \n /abc //notComment { } ~/omg/gg.json", // comment
  /* haha */ "a": 1,
  "b */":/* omg */ "ok",

  "c": true, /* abc
  sfdsfsfs
  */ "arr": [
    "a", // this is a
    "b",
    "c" /*
	sfsd
	*/
  ]
  "sf":{
	"d":"k",  // haha
	"e":"o"
  }  // hahaha
} // omg
`

func Test_UnmarshalJSONC(t *testing.T) {
	r, err := JSONCToJSON([]byte(totalsrc))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(r))
}

func Test_FindSecondLast(t *testing.T) {
	a, b, err := findSecondLastLine([]byte(totalsrc))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(a, b)
}

func Test_Index(t *testing.T) {
	s := "abcbcbc"
	t.Log(strings.Index(s, "d"))
	t.Log(strings.LastIndex(s, "d"))

	t.Log(strings.Index(s, "bc"))
	t.Log(strings.LastIndex(s, "bc"))
}
