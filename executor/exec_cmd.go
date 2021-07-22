package executor

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/mehmetron/pid1/userPort"
	"log"
	"os"
	"os/exec"
	"time"
)

//type Command struct {
//	Operation  string `json:"operation"`
//	Entrypoint string `json:"entrypoint"`
//	Content    string `json:"content"`
//}

type Command struct {
	Operation  string `json:"operation"`
	Entrypoint string `json:"entrypoint"`
	Files      []struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	} `json:"files"`
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

func MultipleFiles(c *Command) {
	fmt.Println("Creating multiple files")

	for _, file := range c.Files {
		f, err := os.Create(file.Name)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		_, err = f.WriteString(file.Content)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//func ListFileNames(c *Command) []string {
//	fmt.Println("Listing file names")
//	//var files []string
//	files := []string{"run"}
//
//	for _, file := range c.Files {
//		files = append(files, file.Name)
//	}
//	return files
//}

func (c *Command) Running() (string, string, uint16, error) {
	demoPort := userPort.GetOpenPortLib()
	packages, err := c.isPackageInstalled()
	fmt.Println(packages)
	fmt.Println("open port: ", demoPort)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	//filemanager.CreateFile(c.Entrypoint, c.Content)
	MultipleFiles(c)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)

	err = InitializeModules()
	if err != nil {
		fmt.Println(err)
	}

	err = InstallPackage("github.com/brianvoe/gofakeit/v6")
	if err != nil {
		fmt.Println(err)
	}

	//files := ListFileNames(c)
	//cmd := exec.CommandContext(ctx, "go", files...)

	// Taken from https://www.reddit.com/r/golang/comments/onjxsn/error_when_trying_to_run_go_project_stat_go_no/
	cmd := exec.CommandContext(ctx, "bash", "-c", "go run *.go")
	//cmd := exec.CommandContext(ctx, "go", "run", "*.go")
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

// ExecStream executes a long executing command and streams the output to a websocket
func (c *Command) ExecStream() {
	cmd := exec.Command("go", "run", "main.go")
	wd, _ := os.Getwd()
	cmd.Dir = fmt.Sprintf("%s/work", wd)

	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}
