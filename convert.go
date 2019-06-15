package slack_lib

import (
	"errors"

	"github.com/nlopes/slack"
)

func ConvertDisplayChannelName(api *slack.Client, ev *slack.MessageEvent) (fromType, name string, err error) {
	// identify channel or group (as known as private channel) or DM

	// Channel prefix : C
	// Group prefix : G
	// Direct message prefix : D
	channel, err := api.GetConversationInfo(ev.Channel, false)
	if err != nil {
		return "", "", err
	}

	name = channel.Name
	if channel.IsChannel && !channel.IsPrivate {
		fromType = "channel"
	} else if channel.IsPrivate {
		fromType = "group"
	} else if channel.IsIM {
		fromType = "DM"
		user, err := api.GetUserInfo(ev.Msg.User)
		if err != nil {
			return fromType, "", nil // NOTE: not error
		}
		name = user.Profile.DisplayName
	} else {
		fromType = "unknown"
	}

	return fromType, name, nil
}

func ConvertDisplayUserName(api *slack.Client, ev *slack.MessageEvent, id string) (username, usertype string, err error) {
	// user id to display name

	if id != "" {
		// specific id (maybe user)
		info, err := api.GetUserInfo(id)
		if err != nil {
			return "", "", err
		}

		return info.Name, "user", nil
	}

	// return self id
	if ev.Msg.BotID == "B01" {
		// this is slackbot
		return "Slack bot", "bot", nil
	} else if ev.Msg.BotID != "" {
		// this is bot
		byInfo, err := api.GetBotInfo(ev.Msg.BotID)
		if err != nil {
			return "", "", err
		}

		return byInfo.Name, "bot", nil
	} else if ev.Msg.SubType != "" {
		// SubType is not define user
		return ev.Msg.SubType, "status", nil
	} else {
		// user
		byInfo, err := api.GetUserInfo(ev.Msg.User)
		if err != nil {
			return "", "", err
		}

		return byInfo.Name, "user", nil
	}

	return "", "", errors.New("user not found")
}