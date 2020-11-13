// +build mage

package main

import (
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

func BuildAll() error {
	err := sh.RunWith(map[string]string{"GOOS": "linux", "GOARCH": "amd64"}, "go", "build", "-o", "./dist/mac_linux_amd64", "./cmd/mac")
	if err != nil {
		return err
	}

	err = sh.RunWith(map[string]string{"GOOS": "darwin", "GOARCH": "amd64"}, "go", "build", "-o", "./dist/mac_darwin_amd64", "./cmd/mac")
	if err != nil {
		return err
	}

	err = sh.RunWith(map[string]string{"GOOS": "windows", "GOARCH": "amd64"}, "go", "build", "-o", "./dist/mac_windows_amd64", "./cmd/mac")
	if err != nil {
		return err
	}

	return nil
}

func Lint() error {
	return sh.RunV("golangci-lint", "run", "./...")
}

func Format() error {
	if err := sh.RunV("gofmt", "-w", "."); err != nil {
		return err
	}

	return nil
}

func Coverage() error {
	// Create a coverage file if it does not already exist
	if err := os.MkdirAll(filepath.Join(".", "coverage"), os.ModePerm); err != nil {
		return err
	}

	if err := sh.RunV("go", "test", "./pkg/...", "-v", "-cover", "-coverprofile=coverage/backend.out"); err != nil {
		return err
	}

	if err := sh.RunV("go", "tool", "cover", "-html=coverage/backend.out", "-o", "coverage/backend.html"); err != nil {
		return err
	}

	return nil
}

func Clean() error {
	err := os.RemoveAll("dist")
	if err != nil {
		return err
	}

	err = os.RemoveAll("coverage")
	if err != nil {
		return err
	}

	return nil
}
