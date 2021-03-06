package neo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"maunium.net/go/mautrix"
	"time"
)

// Event is a wrapper over the mautrix.Event type to permit use without needing to import mautrix
type Event *mautrix.Event

// MatchFunc is used to detect whether the bot should respond to an event
type MatchFunc func(event Event) bool

// HandlerFunc is executed when the corresponding MatchFunc returned true
type HandlerFunc func(c *Context) error

type handler struct {
	Match   MatchFunc
	Handler HandlerFunc
}

// A Bot is the main component of the framework. This struct contains all neccesary informations to interact with a room
type Bot struct {
	logger   *zap.Logger
	client   *mautrix.Client
	handlers map[uuid.UUID]*handler
	interval time.Duration
}

// Configuration is used to initialy configure the bot. Note: changing the configuration after creation will have _no_ effect
type Configuration struct {
	HomeServer      string
	UserID          string
	AccessToken     string
	Logger          *zap.Logger
	PollingInterval time.Duration
}

// NewBot constructs a new Bot with the gien configuration
func NewBot(c *Configuration) (*Bot, error) {
	if c == nil || c.HomeServer == "" || c.AccessToken == "" || c.UserID == "" {
		return nil, errors.New("Configuration must at least include: HomeServer, AccessToken and UserID")
	}
	var logger *zap.Logger
	if c.Logger == nil {
		lg, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		logger = lg
		logger.Info("No logger specified in configuation, defaulting to debug logger")
	}
	client, err := mautrix.NewClient(c.HomeServer, c.UserID, c.AccessToken)
	if err != nil {
		return nil, err
	}

	if c.PollingInterval == 0 {
		c.PollingInterval = time.Second * 10
	}

	logger.Info("initialized bot", zap.String("HomeServer", c.HomeServer), zap.String("UserID", c.UserID))
	return &Bot{
		client:   client,
		handlers: make(map[uuid.UUID]*handler),
		logger:   logger,
		interval: c.PollingInterval,
	}, nil
}

// React is used to route specific events to their respecitve handlers
func (b *Bot) React(m MatchFunc, h HandlerFunc) {
	hid := uuid.New()
	b.handlers[hid] = &handler{Match: m, Handler: h}
	b.logger.Info("registered handler", zap.String("handler_id", hid.String()))
}

// HandleEvent enumerates all handlers and fires the ones matching
func (b *Bot) HandleEvent(e *mautrix.Event) error {
	for hid, h := range b.handlers {
		if h.Match(Event(e)) {
			b.logger.Info("firing handler", zap.String("handler_id", hid.String()))
			err := h.Handler(b.genCtx(e))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Context returns a generic Context without a corresponding event.
// This can be used to independently send messages without needing to react to an Event
func (b *Bot) Context(room string) *Context {
	return b.genCtx(&mautrix.Event{
		RoomID: room,
	})
}

// Run listens for new events on the specified channel and is responsible for dispatching events to the handlers
func (b *Bot) Run(ctx context.Context, room string) error {
	s, err := b.client.SyncRequest(64, "", "", false, "")
	if err != nil {
		b.logger.Error("could not do inital sync", zap.Error(err))
		return err
	}

	name, err := b.client.GetOwnDisplayName()
	if err != nil {
		return err
	}
	b.client.SendText(room, fmt.Sprintf("Hello, I'm %s and now ready for service 🤖", name.DisplayName))

	errors := make(chan error)
	go func() {
		from := s.NextBatch
		for {
			msgs, err := b.client.Messages(room, from, "", 'f', 10)
			if len(msgs.Chunk) == 0 {
				continue
			}
			b.logger.Info("fetched new messages", zap.Int("messages_count", len(msgs.Chunk)))
			if err != nil {
				panic(err)
			}
			for _, event := range msgs.Chunk {
				if event.Sender == b.client.UserID {
					continue
				}
				err := b.HandleEvent(event)
				if err != nil {
					errors <- err
				}
			}
			from = msgs.End
			time.Sleep(b.interval)
		}
	}()
	for {
		select {
		case err = <-errors:
			return err
		case _ = <-ctx.Done():
			return nil
		}
	}
}
