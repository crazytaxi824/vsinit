package util

import (
	"bytes"
	"os/exec"
	"strings"
)

var extSuggestion = Suggestion{
	Problem: "need to install vscode extension ",
	Solution: "you can install it in the vscode extentsion market, or run:\n" +
		"code --install-extension ",
}

// 检查 vscode 和 vscode 插件
func CheckVscodeAndExtensions(exts []string) ([]*Suggestion, error) {
	var suggs []*Suggestion

	sug := CheckCMDInstall("code")
	if sug != nil {
		// vscode 不存在
		suggs = append(suggs, sug, &Suggestion{
			Problem:  extSuggestion.Problem + "'" + strings.Join(exts, ", ") + "'",
			Solution: extSuggestion.Solution + strings.Join(exts, " "),
		})
	} else {
		// vscode 存在
		cmd := exec.Command("code", "--list-extensions")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		var result []string

		for _, ext := range exts {
			if !bytes.Contains(out, []byte(ext)) {
				result = append(result, ext)
			}
		}

		if len(result) > 0 {
			suggs = append(suggs, &Suggestion{
				Problem:  extSuggestion.Problem + "'" + strings.Join(result, ", ") + "'",
				Solution: extSuggestion.Solution + strings.Join(result, " "),
			})
		}
	}

	if len(suggs) > 0 {
		return suggs, nil
	}

	return nil, nil
}
