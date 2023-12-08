package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ChooseSingleFile(fileContents []FileContent, msg string) ([]FileContent, error) {
	fmt.Println("file list:")

	for i, fc := range fileContents {
		fmt.Println(" ", i+1, fc.Filepath)
	}

	total := len(fileContents)
	fmt.Printf("Choose files %s"+msg+"%s: [q(uit)|1~%d]: ", COLOR_BOLD_YELLOW, COLOR_RESET, total)

	// get user input
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("%s: %w", Cancel(), err)
	}
	input = strings.ToLower(strings.TrimSpace(input))

	// cases
	if input == "q" || input == "quit" {
		return nil, fmt.Errorf("%s: %s", Cancel(), "quit")
	}

	// 判断文件编号
	fnum, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", Cancel(), "please choose a valid file number")
	}

	if fnum >= 1 && fnum <= total {
		return []FileContent{fileContents[fnum-1]}, nil
	}

	return nil, fmt.Errorf("%s: %s", Cancel(), "please choose a valid file number")
}
