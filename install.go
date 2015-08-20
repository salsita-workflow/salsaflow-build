package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("Error:", err)
	}
}

func run() error {
	cmdFlag := flag.String("cmd", "go $ install", "command to use to install individual packages")
	flag.Parse()

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

		fmt.Printf("====> Rewriting %v\n", theirFile)

		if err := os.Remove(their); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}

		if err := os.Symlink(our, their); err != nil {
			return err
		}
	}

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

	// Generate environment for running `go install`.
	// Inherit PATH and all GO* variables, except GOPATH.
	gopath := fmt.Sprintf("%v:%v:%v",
		salsaflowWorkspace, salsaflowGodepWorkspace, modulesGodepWorkspace)

	currentEnv := os.Environ()
	env := make([]string, 2, len(currentEnv)+2)
	env[0] = "PATH=" + os.Getenv("PATH")
	env[1] = "GOPATH=" + gopath
	for _, kv := range currentEnv {
		if strings.HasPrefix(kv, "GO") && !strings.HasPrefix(kv, "GOPATH=") {
			env = append(env, kv)
		}
	}

	// Process the command string.
	args := make([]string, 0, 2)
	for _, part := range strings.Split(*cmdFlag, "$") {
		args = append(args, strings.TrimSpace(part))
	}

	// Run the install command for every executable package to be installed.
	packages := []string{
		"github.com/salsaflow/salsaflow",
		"github.com/salsaflow/salsaflow/bin/hooks/salsaflow-commit-msg",
		"github.com/salsaflow/salsaflow/bin/hooks/salsaflow-post-checkout",
		"github.com/salsaflow/salsaflow/bin/hooks/salsaflow-pre-push",
	}

	for _, pkg := range packages {
		fmt.Printf("\n====> Installing %v\n", pkg)

		argv := append(args, pkg)
		fmt.Printf("      cmd = %v\n\n", argv)

		cmd := exec.Command(argv[0], argv[1:]...)
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
