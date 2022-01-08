package util

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func NpmInstallPrefixDev(prefix string, packages []string) error {
	args := []string{"i", "--prefix", prefix, "-D"}
	args = append(args, packages...)

	fmt.Printf("\n%snpm %s%s%s\n", COLOR_YELLOW, COLOR_GREEN, strings.Join(args, " "), COLOR_RESET)

	cmd := exec.Command("npm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NpmInstallDevLocal(packages []string) error {
	args := []string{"i", "-D"}
	args = append(args, packages...)

	fmt.Printf("\n%snpm %s%s%s\n", COLOR_YELLOW, COLOR_GREEN, strings.Join(args, " "), COLOR_RESET)

	cmd := exec.Command("npm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// install react commonly uesd packages
func NpmInstallSaveLocal(packages []string) error {
	args := []string{"i", "-S"}
	args = append(args, packages...)

	fmt.Printf("\n%snpm %s%s%s\n", COLOR_YELLOW, COLOR_GREEN, strings.Join(args, " "), COLOR_RESET)

	cmd := exec.Command("npm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
