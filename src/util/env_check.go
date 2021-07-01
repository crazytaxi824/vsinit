package util

import (
	"bytes"
	"os/exec"
	"strings"
)

var extSuggestion = Suggestion{
	Problem:  "need to install following vscode extensions:",
	Solution: "you can install it in the vscode extentsion market, or run:\n",
}

// 检查 vscode 和 vscode 插件
func CheckVscodeAndExtensions(exts []string) ([]*Suggestion, error) {
	var suggs []*Suggestion

	sug := CheckCMDInstall("code")
	if sug != nil {
		// vscode 不存在
		suggs = append(suggs, sug, &Suggestion{
			Problem: extSuggestion.Problem,
			Solution: extSuggestion.Solution + "code --install-extension " +
				strings.Join(exts, "; \\\ncode --install-extension "),
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
				result = append(result, "code --install-extension "+ext)
			}
		}

		if len(result) > 0 {
			suggs = append(suggs, &Suggestion{
				Problem:  extSuggestion.Problem,
				Solution: extSuggestion.Solution + strings.Join(result, "; \\\n"),
			})
		}
	}

	if len(suggs) > 0 {
		return suggs, nil
	}

	return nil, nil
}
