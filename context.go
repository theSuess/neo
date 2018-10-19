package neo

import (
	"go.uber.org/zap"
	"maunium.net/go/mautrix"
)

type Context struct {
	Logger *zap.Logger
	Event  Event
	client *mautrix.Client
}

func (b *Bot) genCtx(e Event) *Context {
	return &Context{
		Event:  e,
		Logger: b.logger,
		client: b.client,
	}
}

func (c *Context) SendText(text string) error {
	_, err := c.client.SendText(c.Event.RoomID, text)
	return err
}

func (c *Context) UploadLink(link string) (string, error) {
	resp, err := c.client.UploadLink(link)
	return resp.ContentURI, err
}

func (c *Context) UserTyping(typing bool) error {
	_, err := c.client.UserTyping(c.Event.RoomID, typing, 64)
	return err
}

func (c *Context) SendImage(body, url string) error {
	_, err := c.client.SendImage(c.Event.RoomID, body, url)
	return err
}
