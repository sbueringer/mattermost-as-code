package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmdMac provides a cobra command
func NewCmdMac() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "mac",
		Short:   "Mattermost as code",
		Example: "",
		RunE: func(c *cobra.Command, args []string) error {
			return fmt.Errorf("subcommand is mandatory")
		},
	}

	cmd.PersistentFlags().String("username", os.Getenv("MATTERMOST_USERNAME"), "Mattermost username, can also be set via MATTERMOST_USERNAME")
	cmd.PersistentFlags().String("password", os.Getenv("MATTERMOST_PASSWORD"), "Mattermost password, can also be set via MATTERMOST_PASSWORD")
	cmd.PersistentFlags().String("url", os.Getenv("MATTERMOST_URL"), "Mattermost url, can also be set via MATTERMOST_URL")

	cmd.AddCommand(NewCmdImport())
	cmd.AddCommand(NewCmdExport())
	return cmd
}
