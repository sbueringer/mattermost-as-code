package main

import (
	"os"

	"github.com/sbueringer/mattermost-as-code/pkg/cmd"
)

func main() {
	if err := cmd.NewCmdMac().Execute(); err != nil {
		os.Exit(1)
	}
}
