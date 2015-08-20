package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func run() error {
	fileNames := []string{
		"issue_tracker_modules.go",
		"code_review_tool_modules.go",
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	salsaflowWorkspace := filepath.Join(cwd, "workspace")

	for _, name := range fileNames {
		var (
			ourFile   = filepath.Join("assets", name)
			our       = filepath.Join(cwd, ourFile)
			theirFile = filepath.Join("workspace/src/github.com/salsaflow/salsaflow/modules", name)
			their     = filepath.Join(cwd, theirFile)
		)

		fmt.Printf("---> Rewriting %v\n", theirFile)

		if err := os.Remove(their); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}

		if err := os.Symlink(our, their); err != nil {
			return err
		}
	}

	fmt.Println()

	// Get Godep workspace for custom modules.
	modulesGodepWorkspace, err := godepWorkspace(cwd)
	if err != nil {
		return err
	}

	// Get Godep workspace for SalsaFlow.
	salsaflowGodepWorkspace, err := godepWorkspace(
		filepath.Join(salsaflowWorkspace, "src/github.com/salsaflow/salsaflow"))
	if err != nil {
		return err
	}

	// Install SalsaFlow.
	var (
		path   = os.Getenv("PATH")
		goroot = os.Getenv("GOROOT")
		gopath = os.Getenv("GOPATH")
	)

	gopath = fmt.Sprintf("%v:%v:%v:%v",
		salsaflowWorkspace, salsaflowGodepWorkspace, modulesGodepWorkspace, gopath)

	env := []string{
		"PATH=" + path,
		"GOROOT=" + goroot,
		"GOPATH=" + gopath,
	}

	packages := []string{
		"github.com/salsaflow/salsaflow",
		"github.com/salsaflow/salsaflow/bin/hooks/salsaflow-commit-msg",
		"github.com/salsaflow/salsaflow/bin/hooks/salsaflow-post-checkout",
		"github.com/salsaflow/salsaflow/bin/hooks/salsaflow-pre-push",
	}

	for _, pkg := range packages {
		fmt.Printf("---> go install %v\n", pkg)

		cmd := exec.Command("go", "install", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = env

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func godepWorkspace(wd string) (string, error) {
	cmd := exec.Command("godep", "path")
	cmd.Dir = wd

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(output)), nil
}
