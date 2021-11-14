package controller

import (
	"chat/v1/pkg/model/service/channel"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChannelResolver interface {
	StreamHandler(*gin.Context)
	SendMassageHandler(*gin.Context)
}

type channelResolver struct {
	manager *channel.ChannelManager
}

func (c *Controller) Channel() ChannelResolver {
	return &channelResolver{manager: c.channel}
}

func (c *channelResolver) StreamHandler(cxt *gin.Context) {
	channel := cxt.Param("id")
	user := cxt.Param("user")
	listener := c.manager.Open(channel)
	defer c.manager.Close(channel, listener)

	c.manager.Post("System", channel,
		fmt.Sprintf("login %s", user),
	)
	cxt.Stream(func(w io.Writer) bool {
		select {
		case message := <-listener:
			cxt.SSEvent("message", message)
		case <-cxt.Writer.CloseNotify():
			return false
		}
		return true
	})
}

func (c *channelResolver) SendMassageHandler(cxt *gin.Context) {
	var req postMessageRequest
	channel := cxt.Param("id")
	user := cxt.GetHeader("user")

	if err := cxt.BindJSON(&req); err != nil {
		cxt.JSON(
			http.StatusInternalServerError,
			gin.H{"status": "Post failed."},
		)
	}
	c.manager.Post(user, channel, req.Content)
	cxt.JSON(
		http.StatusOK,
		gin.H{"status": "Post success."},
	)
}
