package channel

import (
	"fmt"
	"time"

	"github.com/dustin/go-broadcast"
)

var messageFormat = `%s %s
> %s`

const (
	CHANNEL_SIZE = 100
)

type listener struct {
	ID string
	Ch chan interface{}
}

type message struct {
	userID    string
	channelID string
	content   string
}

type ChannelManager struct {
	channels map[string]broadcast.Broadcaster

	open  chan *listener
	close chan *listener

	delete chan string

	messages chan *message
}

func NewChannelManager() ChannelManager {
	manager := ChannelManager{
		channels: make(map[string]broadcast.Broadcaster),
		open:     make(chan *listener, CHANNEL_SIZE),
		close:    make(chan *listener, CHANNEL_SIZE),
		delete:   make(chan string, CHANNEL_SIZE),
		messages: make(chan *message, CHANNEL_SIZE),
	}

	go manager.run()
	return manager
}

func (c *ChannelManager) run() {
	for {
		select {
		case listener := <-c.open:
			c.channel(listener.ID).Register(listener.Ch)
		case listener := <-c.close:
			c.channel(listener.ID).Unregister(listener.Ch)
			close(listener.Ch)
		case id := <-c.delete:
			broadcaster, ok := c.channels[id]
			if ok {
				broadcaster.Close()
				delete(c.channels, id)
			}
		case message := <-c.messages:
			content := fmt.Sprintf(messageFormat,
				message.userID,
				time.Now().String(),
				message.content,
			)
			c.channel(message.channelID).Submit(
				content,
			)
		}
	}
}

func (c *ChannelManager) channel(id string) broadcast.Broadcaster {
	broadcaster, ok := c.channels[id]
	if !ok {
		broadcaster = broadcast.NewBroadcaster(10)
		c.channels[id] = broadcaster
	}
	return broadcaster
}

func (c *ChannelManager) Open(id string) chan interface{} {
	lis := make(chan interface{})
	c.open <- &listener{
		ID: id,
		Ch: lis,
	}
	return lis
}

func (c *ChannelManager) Close(id string, lis chan interface{}) {
	c.close <- &listener{
		ID: id,
		Ch: lis,
	}
}

func (c *ChannelManager) Post(userID, channelID, content string) {
	msg := &message{
		userID:    userID,
		channelID: channelID,
		content:   content,
	}
	c.messages <- msg
}

func (c *ChannelManager) Delete(id string) {
	c.delete <- id
}
