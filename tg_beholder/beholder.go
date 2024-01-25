package tgcrawl

import (
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
)


func (tgc *TgCrawler) CheckedPosts() error{
	api := tgc.client.API()

	d := tg.NewUpdateDispatcher()

	gaps := updates.New(
		updates.Config{
			Handler: d,
			Logger: tgc.Logger,
		})


	msgs, err := api.

	https://github.com/gotd/td/blob/main/examples/updates/main.go

	return nil
}