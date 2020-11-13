package cmd

import (
	"strings"

	"github.com/spf13/pflag"

	"github.com/mattermost/mattermost-server/v5/model"
)

func getClient(flags *pflag.FlagSet) (*model.Client4, *model.User, error) {
	username, err := flags.GetString("username")
	if err != nil {
		return nil, nil, err
	}
	password, err := flags.GetString("password")
	if err != nil {
		return nil, nil, err
	}
	url, err := flags.GetString("url")
	if err != nil {
		return nil, nil, err
	}

	client := model.NewAPIv4Client(url)

	user, resp := client.Login(username, password)
	if resp.Error != nil {
		return nil, nil, err
	}
	return client, user, nil
}

func getChannelsForTeamForUserByID(client *model.Client4, teamID, userID string) (map[string]*model.Channel, error) {
	channels, resp := client.GetChannelsForTeamForUser(teamID, userID, false, "")
	if resp.Error != nil {
		return nil, resp.Error
	}
	ret := map[string]*model.Channel{}
	for _, channel := range channels {
		ret[channel.Id] = channel
	}
	return ret, nil
}

func getChannelsForTeamForUserByName(client *model.Client4, teamID, userID string) (map[string]*model.Channel, error) {
	channels, resp := client.GetChannelsForTeamForUser(teamID, userID, false, "")
	if resp.Error != nil {
		return nil, resp.Error
	}
	ret := map[string]*model.Channel{}
	for _, channel := range channels {
		channelName, err := calculateChannelName(client, userID, channel)
		if err != nil {
			return nil, err
		}
		ret[channelName] = channel
	}
	return ret, nil
}

func calculateChannelName(client *model.Client4, userID string, ch *model.Channel) (string, error) {
	if ch.Type == model.CHANNEL_DIRECT {
		members, resp := client.GetChannelMembers(ch.Id, 0, 100, "")
		if resp.Error != nil {
			return "", resp.Error
		}
		var membersArr []string
		for _, member := range *members {
			if member.UserId == userID {
				continue
			}
			user, resp := client.GetUser(member.UserId, "")
			if resp.Error != nil {
				return "", resp.Error
			}
			membersArr = append(membersArr, user.Username)
		}
		return strings.Join(membersArr, " "), nil
	}
	if ch.Type == model.CHANNEL_GROUP {
		return ch.DisplayName, nil
	}
	return ch.Name, nil
}
