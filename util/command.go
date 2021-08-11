package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

// PatternExists returns true if the given glob matches any file in
// the current directory.
func PatternExists(pattern string) bool {
	if matches, err := filepath.Glob(pattern); err != nil {
		panic(err)
	} else {
		return len(matches) > 0
	}
}

// Exists returns true if a directory entry by the given filename
// exists. If an I/O error occurs, FileExists terminates the process.
func Exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else if err != nil {
		fmt.Printf("%s: %s\n", filename, err)
		return false
	} else {
		return true
	}
}

// TempDir creates and returns the name of temporary directory. If
// creation fails, it terminates the process. The caller is
// responsible for cleaning up the temporary directory afterwards.
func TempDir() string {
	if tempdir, err := ioutil.TempDir("", "upm"); err != nil {
		fmt.Println(err)
		return ""
	} else {
		return tempdir
	}
}

// Regexps compiles each provided pattern into a regexp object, and
// returns a slice of them.
func Regexps(patterns []string) []*regexp.Regexp {
	regexps := []*regexp.Regexp{}
	for _, pattern := range patterns {
		regexps = append(regexps, regexp.MustCompile(pattern))
	}
	return regexps
}

// RunCmd prints and runs the given command, exiting the process on
// error or command failure. Stdout and stderr go to the terminal.
func RunCmd(cmd []string) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stderr
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		fmt.Println("66  ", err)
	}
}

// GetCmdOutput prints and runs the given command, returning its
// stdout as a string. Stderr goes to the terminal. GetCmdOutput exits
// the process on error or command failure.
func GetCmdOutput(cmd []string) []byte {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stderr = os.Stderr
	output, err := command.Output()
	if err != nil {
		fmt.Println("78 ", err)
	}
	return output
}

// GetExitCode runs a commands, and optionally prints the output to
// stdout and/or stderr, and it returns the exit code afterwards.
func GetExitCode(cmd []string, printStdout bool, printStderr bool) int {
	command := exec.Command(cmd[0], cmd[1:]...)
	if printStdout {
		command.Stdout = os.Stdout
	}
	if printStderr {
		command.Stderr = os.Stderr
	}
	if err := command.Run(); err != nil {
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}
