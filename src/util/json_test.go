package util

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func Test_findLastChar(t *testing.T) {
	src := `"s": "abc\" def\" \\\t \n /abc //notComment { } ~/omg/gg.json", // comment`

	lastIndex, multiLineComment, err := lastValidCharInJSONCline([]byte(src), 0)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(src[:lastIndex+1])
	t.Log(multiLineComment)
}

func Test_findLastChar2(t *testing.T) {
	src := `{`

	lastIndex, multiLineComment, err := lastValidCharInJSONCline([]byte(src), 0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(lastIndex)
	t.Log(src[:lastIndex+1])
	t.Log(multiLineComment)
}

func Test_findLastChar3(t *testing.T) {
	src := `{ "a": "c" }`

	lastIndex, multiLineComment, err := lastValidCharInJSONCline([]byte(src), 0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(lastIndex)
	t.Log(src[:lastIndex+1])
	t.Log(multiLineComment)
}

const totalsrc = `{
  // this is comment
  "s": "abc\" def\" \\\t \n /abc //notComment { } ~/omg/gg.json", // comment
  /* haha */ "a": /*dfsd*/ 1,
  "b */":/* omg */ "ok",
  /*omg
  hahaha
  */

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

func Test_JSONCToJSON(t *testing.T) {
	r, err := JSONCToJSON([]byte(totalsrc))
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

func Test_AppendJSONC(t *testing.T) {
	n, err := AppendToJSONC([]byte(totalsrc), []byte(`  "o":1,
  "k":2`))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(n))
}

func Test_AppendJSONC2(t *testing.T) {
	n, err := AppendToJSONC([]byte("{}"), []byte(`  "o":1,
  "k":2`))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(n))
}

func Test_AppendJSONC3(t *testing.T) {
	n, err := AppendToJSONC([]byte(`{"a":"b"}`), []byte(`  "o":1,
  "k":2`))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(n))
}
func Test_Index(t *testing.T) {
	s := "abcbcbc"
	t.Log(strings.Index(s, "d"))
	t.Log(strings.LastIndex(s, "d"))

	t.Log(strings.Index(s, "bc"))
	t.Log(strings.LastIndex(s, "bc"))
}
