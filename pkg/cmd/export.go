package cmd

import (
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/spf13/cobra"
)

var (
	team string
)

// NewCmdExport provides the export command
func NewCmdExport() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export",
		Short:   "Export the channel configuration",
		Example: "",
		RunE: func(c *cobra.Command, _ []string) error {
			return runExport(c)
		},
	}
	cmd.Flags().StringVar(&team, "input", "", "export only a specific team")
	return cmd
}

func runExport(c *cobra.Command) error {
	client, user, err := getClient(c.Flags())
	if err != nil {
		return err
	}

	teams, resp := client.GetTeamsForUser(user.Id, "")
	if resp.Error != nil {
		return err
	}

	outputTeams := Teams{}
	for _, team := range teams {

		outputTeam := Team{Name: team.Name}

		channelsMapping, err := getChannelsForTeamForUserByID(client, team.Id, user.Id)
		if err != nil {
			return err
		}

		categories, resp := client.GetSidebarCategoriesForTeamForUser(user.Id, team.Id, "")
		if resp.Error != nil {
			return err
		}

		for _, category := range categories.Categories {
			// don't export the direct_messages and channels categories
			// as they are both just the not otherwise categorized channels
			if category.Type == model.SidebarCategoryDirectMessages || category.Type == model.SidebarCategoryChannels {
				continue
			}

			outputCategory := Category{
				Name:    category.DisplayName,
				Type:    string(category.Type),
				Sorting: string(category.Sorting),
			}

			for _, channelID := range category.Channels {
				ch, ok := channelsMapping[channelID]
				if !ok {
					return fmt.Errorf("channel %s not found", channelID)
				}
				channelName, err := calculateChannelName(client, user.Id, ch)
				if err != nil {
					return err
				}
				outputCategory.Channels = append(outputCategory.Channels, channelName)
			}

			outputTeam.Categories = append(outputTeam.Categories, outputCategory)
		}

		outputTeams.Teams = append(outputTeams.Teams, outputTeam)
	}

	output, err := yaml.Marshal(outputTeams)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
