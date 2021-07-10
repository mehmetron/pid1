package main

import "os/exec"

type PlayError struct {
	err    error
	output string
}

func (e *PlayError) Error() string {
	return e.err.Error()
}

func (e *PlayError) Output() string {
	return e.output
}

func Goimports() (string, error) {
	cmd := exec.Command("goimports", "-w", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &PlayError{
			err:    err,
			output: string(output),
		}
	}
	return string(output), nil
}

func PlayGoFile(file string) (string, error) {
	cmd := exec.Command("go", "run", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &PlayError{
			err:    err,
			output: string(output),
		}
	}

	return string(output), nil
}
