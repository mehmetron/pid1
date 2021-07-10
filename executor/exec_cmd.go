package executor

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mehmetron/pid1/filemanager"
	"github.com/mehmetron/pid1/userPort"
	"log"
	"os"
	"os/exec"
	"time"
)

type Command struct {
	Operation  string `json:"operation"`
	Entrypoint string `json:"entrypoint"`
	Content    string `json:"content"`
}

func Handler() string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", "ls && sleep 5 && cd .. && ls")
	//cmd.Dir = "/Users/mehmetcureoglu/go/src"

	var stdoutBuf, stderrBuf bytes.Buffer

	// 1
	//cmd.Stdout = io.MultiWriter(w, os.Stdout, &stdoutBuf)
	//cmd.Stderr = io.MultiWriter(w, os.Stderr, &stderrBuf)

	// 2
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// 3
	//pr, pw := io.Pipe()
	//defer pw.Close()
	//cmd.Stdout = pw
	//cmd.Stderr = pw
	//go io.Copy(w, pr)

	// We ignore the returned error because
	// If the command was killed, err will be "signal: killed"
	// If the command wasn't killed, it contains the actual error, e.g. invalid command
	cmd.Run()
	//err := cmd.Run()
	//if err != nil {
	//	log.Fatalf("cmd.Run() failed with %s\n", err)
	//}

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command was killed")
		//fmt.Fprintf(w, "Command was killed")
		return ""
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

	// 2
	return fmt.Sprintf("Output: %s \n Err: %s", outStr, errStr)
	//outStr, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//if (ctx.Err() == context.DeadlineExceeded) {
	//	fmt.Println("Command was killed")
	//	fmt.Fprintf(w, "Command was killed")
	//	return
	//}
	//fmt.Println(fmt.Sprintf("Output: %s \n", outStr))
	//fmt.Fprintf(w, fmt.Sprintf("Output: %s \n", outStr))
}

func PipingStd() {
	fmt.Println("Hello, playground")
	cmd := exec.Command("ls", "-lah")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func (c *Command) Running() (string, string, uint16, error) {
	demoPort := userPort.GetOpenPortLib()
	packages, err := c.isPackageInstalled()
	fmt.Println(packages)
	fmt.Println("open port: ", demoPort)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	filemanager.CreateFile(c.Entrypoint, c.Content)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", c.Entrypoint)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	return stdout.String(), stderr.String(), demoPort, nil
}

func (c *Command) isPackageInstalled() ([]byte, error) {
	fmt.Println("Check if package installed")

	return []byte{'g', 'o', 'l', 'a', 'n', 'g'}, nil
}
