//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// execProdivedCmd executes a command along with its env variables
func execProdivedCmd(tailArgs []string, chdirPath string) (err error) {
	ps, err := exec.LookPath("powershell.exe")
	if err != nil {
		return fmt.Errorf("executable \"powershell.exe\" was not found\n%s", err)
	}
	args := []string{"-NoProfile", "-NonInteractive", "-Command"}
	args = append(args, "$ErrorActionPreference = \"Stop\"; ")
	args = append(args, tailArgs...)
	args = append(args, "; if ($LastExitCode -gt 0) { exit $LastExitCode };")
	cmd := exec.Command(ps, args...)
	cmd.Dir = chdirPath
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
