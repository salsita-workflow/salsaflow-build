package main

import (
	"fmt"
	"log"
	"os"
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

	for _, name := range fileNames {
		var (
			ourFile   = name + ".template"
			our       = filepath.Join(cwd, ourFile)
			theirFile = filepath.Join("salsaflow/modules", name)
			their     = filepath.Join(cwd, theirFile)
		)

		fmt.Printf("---> Removing %v\n", theirFile)
		if err := os.Remove(their); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}

		fmt.Printf("---> Linking %v to act as %v\n", ourFile, theirFile)
		if err := os.Symlink(our, their); err != nil {
			return err
		}

		fmt.Println()
	}

	return nil
}
