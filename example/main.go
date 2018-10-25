package main

import (
	"context"
	"fmt"
	"github.com/theSuess/neo"
	"go.uber.org/zap"
	"os"
	"regexp"
)

func main() {
	username := os.Getenv("BOT_NAME")
	accesskey := os.Getenv("BOT_TOKEN")
	homeserver := "https://chat.weho.st"
	b, err := neo.NewBot(&neo.Configuration{
		HomeServer:  homeserver,
		UserID:      username,
		AccessToken: accesskey,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	room := "!KSkuPboAFRAJICRXFr:matrix.org"
	reg, _ := regexp.Compile(`(?i)you`)
	b.React(
		func(e neo.Event) bool {
			return reg.MatchString(e.Content.Body)
		},
		func(ctx *neo.Context) error {
			ctx.Logger.Info("recieved message", zap.String("content", ctx.Event.Content.Body), zap.Any("event", ctx.Event))
			return ctx.SendText("Did you mean: \"" + reg.ReplaceAllString(ctx.Event.Content.Body, "we") + "\"")
		})
	b.React(
		func(e neo.Event) bool {
			return true
		},
		func(ctx *neo.Context) error {
			ctx.Logger.Info("recieved message", zap.String("content", ctx.Event.Content.Body), zap.Any("event", ctx.Event))
			return nil
		})
	b.Run(context.Background(), room)
}
