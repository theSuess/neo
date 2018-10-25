package neo

import (
	"go.uber.org/zap"
	"maunium.net/go/mautrix"
)

// Context wraps important metadata and allows for the creation of new messages and events
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

// SendText sends the text to the current room
func (c *Context) SendText(text string) error {
	_, err := c.client.SendText(c.Event.RoomID, text)
	return err
}

// UploadLink uploads a remote file to the matrix server. This method returns a matrix url which can then be used
// to send images or other media
func (c *Context) UploadLink(link string) (string, error) {
	resp, err := c.client.UploadLink(link)
	return resp.ContentURI, err
}

// UserTyping indicates user activity on the bots part
func (c *Context) UserTyping(typing bool) error {
	_, err := c.client.UserTyping(c.Event.RoomID, typing, 64)
	return err
}

// SendImage sends a new message with the contents of body and an attached image specified by url.
// This url has to be a matrix content reference
func (c *Context) SendImage(body, url string) error {
	_, err := c.client.SendImage(c.Event.RoomID, body, url)
	return err
}
