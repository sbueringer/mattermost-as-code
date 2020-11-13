package cmd

import (
	"fmt"

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

	cmd.PersistentFlags().String("username", "", "Mattermost username, can also be set via MATTERMOST_USERNAME")
	cmd.PersistentFlags().String("password", "", "Mattermost password, can also be set via MATTERMOST_PASSWORD")
	cmd.PersistentFlags().String("url", "", "Mattermost url, can also be set via MATTERMOST_URL")

	cmd.AddCommand(NewCmdImport())
	cmd.AddCommand(NewCmdExport())
	return cmd
}
