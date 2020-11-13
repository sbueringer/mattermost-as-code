package cmd

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/ghodss/yaml"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/spf13/cobra"
)

var (
	input string
)

// NewCmdExport provides the import command
func NewCmdImport() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "import",
		Short:   "Mattermost as code",
		Example: "",
		RunE: func(c *cobra.Command, _ []string) error {
			return runImport(c)
		},
	}
	cmd.Flags().StringVar(&input, "input", "", "input file")
	return cmd
}

func runImport(c *cobra.Command) error {

	inputTeams, err := readTeams()
	if err != nil {
		return err
	}

	client, user, err := getClient(c.Flags())
	if err != nil {
		return err
	}

	for _, inputTeam := range inputTeams.Teams {
		team, resp := client.GetTeamByName(inputTeam.Name, "")
		if resp.Error != nil {
			return resp.Error
		}

		channelsMapping, err := getChannelsForTeamForUserByName(client, team.Id, user.Id)
		if err != nil {
			return err
		}

		favCategoryId, channelsCategoryId, directMessagesId, err := getSpecialCategoriesIDs(client, user.Id, team.Id)
		if err != nil {
			return err
		}

		if err := removeExistingCustomCategories(client, user.Id, team.Id); err != nil {
			return err
		}

		var categoryOrder []string
		var favCategoryAdded, channelsCategoryAdded, directMessagesCategoryAdded bool
		for _, inputCategory := range inputTeam.Categories {

			// get channel id based on channel name
			var channelIDs []string
			for _, channelName := range inputCategory.Channels {
				ch, ok := channelsMapping[channelName]
				if !ok {
					return fmt.Errorf("channel %s not found", channelName)
				}
				channelIDs = append(channelIDs, ch.Id)
			}

			category := &model.SidebarCategoryWithChannels{
				SidebarCategory: model.SidebarCategory{
					UserId:      user.Id,
					TeamId:      team.Id,
					Sorting:     model.SidebarCategorySorting(inputCategory.Sorting),
					Type:        model.SidebarCategoryType(inputCategory.Type),
					DisplayName: inputCategory.Name,
				},
				Channels: channelIDs,
			}

			// create custom sidebar category
			if inputCategory.Type == string(model.SidebarCategoryCustom) {
				out, resp := client.CreateSidebarCategoryForTeamForUser(user.Id, team.Id, category)
				if resp.Error != nil {
					return resp.Error
				}
				categoryOrder = append(categoryOrder, out.Id)
			}

			// update favorite sidebar category
			if inputCategory.Type == string(model.SidebarCategoryFavorites) {
				_, resp = client.UpdateSidebarCategoryForTeamForUser(user.Id, team.Id, favCategoryId, category)
				if resp.Error != nil {
					return resp.Error
				}
				categoryOrder = append(categoryOrder, favCategoryId)
				favCategoryAdded = true
			}

			// just add it to the order if it was mentioned in the config file
			if inputCategory.Type == string(model.SidebarCategoryChannels) {
				categoryOrder = append(categoryOrder, channelsCategoryId)
				channelsCategoryAdded = true
			}

			// just add it to the order if it was mentioned in the config file
			if inputCategory.Type == string(model.SidebarCategoryDirectMessages) {
				categoryOrder = append(categoryOrder, directMessagesId)
				directMessagesCategoryAdded = true
			}
		}

		// add the default categories to the end of the order if they
		// where not explicitly configured
		if !favCategoryAdded {
			categoryOrder = append(categoryOrder, favCategoryId)
		}
		if !channelsCategoryAdded {
			categoryOrder = append(categoryOrder, channelsCategoryId)
		}
		if !directMessagesCategoryAdded {
			categoryOrder = append(categoryOrder, directMessagesId)
		}

		_, resp = client.UpdateSidebarCategoryOrderForTeamForUser(user.Id, team.Id, categoryOrder)
		if resp.Error != nil {
			return resp.Error
		}
	}

	return nil
}

func readTeams() (*Teams, error) {
	teams := Teams{}

	inputFile, err := ioutil.ReadFile(input)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(inputFile, &teams); err != nil {
		return nil, err
	}

	return &teams, nil
}

func getSpecialCategoriesIDs(client *model.Client4, userID, teamID string) (string, string, string, error) {
	cs, resp := client.GetSidebarCategoriesForTeamForUser(userID, teamID, "")
	if resp.Error != nil {
		return "", "", "", resp.Error
	}

	var favCategoryID string
	var channelsCategoryID string
	var directMessagesID string
	for _, c := range cs.Categories {
		if c.Type == model.SidebarCategoryFavorites {
			favCategoryID = c.Id
		}
		if c.Type == model.SidebarCategoryChannels {
			channelsCategoryID = c.Id
		}
		if c.Type == model.SidebarCategoryDirectMessages {
			directMessagesID = c.Id
		}
	}
	if favCategoryID == "" {
		return "", "", "", fmt.Errorf("could not find favorite category")
	}
	if channelsCategoryID == "" {
		return "", "", "", fmt.Errorf("could not find channels category")
	}
	if directMessagesID == "" {
		return "", "", "", fmt.Errorf("could not find direct messages category")
	}
	return favCategoryID, channelsCategoryID, directMessagesID, nil
}

func removeExistingCustomCategories(client *model.Client4, userID, teamID string) error {
	oldCategories, resp := client.GetSidebarCategoriesForTeamForUser(userID, teamID, "")
	if resp.Error != nil {
		return resp.Error
	}
	for _, oldCategory := range oldCategories.Categories {
		if oldCategory.Type == model.SidebarCategoryCustom {
			route := client.GetUserCategoryRoute(userID, teamID)
			if _, appErr := client.DoApiDelete(path.Join(route, oldCategory.Id)); appErr != nil {
				return appErr
			}
		}
	}
	return nil
}
