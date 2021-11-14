package controller

import (
	"chat/v1/pkg/model/service/channel"
)

type Controller struct {
	channel *channel.ChannelManager
}

func NewController() Controller {
	channelManager := channel.NewChannelManager()
	return Controller{
		channel: &channelManager,
	}
}
