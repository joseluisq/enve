//go:build windows
// +build windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// execProdivedCmd executes a command along with its env variables
func execProdivedCmd(tailArgs []string, chdirPath string, newEnv bool, envVars []string) (err error) {
	ps, err := exec.LookPath("powershell.exe")
	if err != nil {
		return fmt.Errorf("error: executable 'powershell.exe' was not found.\n%v", err)
	}
	args := []string{"-NoProfile", "-NonInteractive", "-Command"}
	args = append(args, "$ErrorActionPreference = \"Stop\"; ")
	args = append(args, tailArgs...)
	args = append(args, "; if ($LastExitCode -gt 0) { exit $LastExitCode };")
	cmd := exec.Command(ps, args...)
	cmd.Dir = chdirPath
	if newEnv {
		cmd.Env = envVars
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
