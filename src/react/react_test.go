package react

import (
	"encoding/json"
	"strings"
	"testing"
)

func Test_DisallowUnknownFields(t *testing.T) {
	type person struct {
		Name string `json:"name,omitempty"`
	}

	var p person
	de := json.NewDecoder(strings.NewReader(`{"name":"kk","age":123}`))
	de.DisallowUnknownFields()
	err := de.Decode(&p)
	if err != nil {
		t.Error(err)
		return
	}
}
