package util

import (
	"bytes"
	"encoding/json"
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
  /* haha */ "a": /*dfsd*/ 1,
  "b */":/* omg */ "ok",

  "c": true, /* abc
  sfdsfsfs
  */ "arr": [
    "a", // this is a
    "b",
    "c" /*
	sfsd
	*/ // dfhskjf*/
  ], /* dfhsf
  fsdfdsf
  */
  "sf":{
	"d":"k",  // haha
	"e":"o"
  }  // hahaha
} // omg

// this is the {comment} after setting
`

func Test_Json2(t *testing.T) {
	r, err := JSONCToJSON2([]byte(totalsrc))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(r))

	var buf bytes.Buffer
	err = json.Indent(&buf, r, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(buf.String())
}

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