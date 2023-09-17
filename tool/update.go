package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func update(args []string) error {
	// update can only be runnable in the project root
	if _, err := os.Stat("go.mod"); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("workdir is not the project root")
	}

	log.Println("Copy $GOROOT/misc/wasm/wasm_exec.js")
	goroot := findGOROOT()
	src := filepath.Join(goroot, "misc", "wasm", "wasm_exec.js")
	dst := "wasm_exec.js"
	if err := copyFile(dst, src); err != nil {
		return fmt.Errorf("copy wasm_exec.js: %w", err)
	}

	log.Println("Update go (go get go)")
	cmd := exec.Command("go", "get", "go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get go: %w", err)
	}

	log.Println("Update dependencies (go get all)")
	cmd = exec.Command("go", "get", "all")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go get all: %w", err)
	}

	log.Println("Cleanup dependencies (go mod tidy)")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	return nil
}

func findGOROOT() string {
	if env := os.Getenv("GOROOT"); env != "" {
		return filepath.Clean(env)
	}
	def := filepath.Clean(runtime.GOROOT())
	out, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return def
	}
	return strings.TrimSpace(string(out))
}
