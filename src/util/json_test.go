package util

import (
	"testing"
)

func Test_FindLastChar(t *testing.T) {
	src := `"s": "abc\" def\" \\\t \n /abc //notComment { } ~/omg/gg.json", // comment`

	lastIndex, hasComments, err := analyseJSONCstatement([]byte(src))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(src[:lastIndex+1])
	t.Log(hasComments)
}

func Test_UnmarshalJSONC(t *testing.T) {
	src := `{
  // this is comment
  "s": "abc\" def\" \\\t \n /abc //notComment { } ~/omg/gg.json", // comment
  "a": 1,
  "b": "ok",


  "c": true,
  "arr": [
    "a", // this is a
    "b",
    "c"
  ]
}
`
	r, err := JSONCToJSON([]byte(src))
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(string(r))
}
