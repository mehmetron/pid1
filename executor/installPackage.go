package executor

import (
	"fmt"
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

func InitializeModules() error {
	defer fmt.Println("Done initializing modules...")
	cmd := exec.Command("go", "mod", "init", "stuff")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
