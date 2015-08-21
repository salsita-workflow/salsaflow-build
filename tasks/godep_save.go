package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var packages = []string{
	"modules/jira",
	"modules/reviewboard",
}

func main() {
	if err := run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("godep", append([]string{"save"}, packages...)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{
		fmt.Sprintf("PATH=%v", os.Getenv("PATH")),
		fmt.Sprintf("GOPATH=%v:%v", os.Getenv("GOPATH"), filepath.Join(cwd, "workspace")),
	}

	if err := cmd.Run(); err != nil {
		log.Fatalln("Error:", err)
	}

	return nil
}
