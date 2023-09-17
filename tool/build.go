package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func build(args []string) error {
	// Parse flags
	flag := flag.NewFlagSet("build", flag.ExitOnError)
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./tool build [arguments]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	addr := flag.String("http", defaultAddr, "HTTP service address")
	flag.Parse(args)

	if flag.NArg() > 0 {
		fmt.Fprintln(os.Stderr, "Unexpected arguments:", flag.Args())
		flag.Usage()
	}

	// Run go build
	cmd := exec.Command("go", "build", "-o", "game.wasm")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}

	// After building, send a request to '_notify' to automatically reload the browser
	u := url.URL{
		Scheme: "http",
		Host:   *addr,
		Path:   "/_notify",
	}

	// Ignore the error, as the build can be done even if the server is not running
	http.PostForm(u.String(), nil)

	return nil
}
