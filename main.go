package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// InstallPackage installs the package passed in
func InstallPackage(pkg string) error {
    defer fmt.Println("Done installing...")
	cmd := exec.Command("go", "get", "-v", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func WatchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/tmp/foo")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {
	fmt.Println("hello world")
	err := InstallPackage("github.com/ermos/gomon")
	if err != nil {
	    fmt.Println(err)
    }


}



