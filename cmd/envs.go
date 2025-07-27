//go:build !windows
// +build !windows

package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// execProdivedCmd executes a command along with its env variables
func execProdivedCmd(tailArgs []string, chdirPath string, newEnv bool, envVars []string) (err error) {
	cmdIn := tailArgs[0]
	c, err := exec.LookPath(cmdIn)
	if err != nil {
		return fmt.Errorf("executable \"%s\" was not found\n%s", cmdIn, err)
	}
	cmd := exec.Command(c, tailArgs[1:]...)
	cmd.Dir = chdirPath
	if newEnv {
		cmd.Env = envVars
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
