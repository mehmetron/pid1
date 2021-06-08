package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)


type Command struct {
	Operation string `json:"operation"`
	Entrypoint string `json:"entrypoint"`
	Content string `json:"content"`
}

func (c *Command) running() (string, string, error){
	port, err := c.isPortOpened()
	packages, err := c.isPackageInstalled()
	fmt.Println(packages)
	fmt.Println(port)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx,"go", "run", c.Entrypoint)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	return stdout.String(), stderr.String(), nil
}

func (c *Command) isPortOpened() (int, error) {
	fmt.Println("Check if port opened")

	return 0, nil
}

func (c *Command) isPackageInstalled() ([]byte, error) {
	fmt.Println("Check if package installed")

	return []byte{'g', 'o', 'l', 'a', 'n', 'g'}, nil
}


