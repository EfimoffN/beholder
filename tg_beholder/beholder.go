package tg_beholder

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/updates"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func (tgc *TgBeholder) CheckedPosts() error {
	log := zap.NewExample()

	api := tgc.client.API()

	self, err := tgc.client.Self(tgc.ctx)
	if err != nil {
		log.Info(err.Error())
		return err
	}

	// d := tg.NewUpdateDispatcher()

	// gaps := updates.New(
	// 	updates.Config{
	// 		Handler: d,
	// 		Logger:  log,
	// 	})

	// Setup message update handlers.
	tgc.dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		if update.Message.String() == "" {
			log.Info("Message", zap.Any("message", update.Message))
		}

		log.Info("Channel message", zap.Any("message", update.Message))
		return nil
	})

	// Create message sending helper.
	sender := message.NewSender(tgc.client.API())

	tgc.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		// if update.Message.String() == "" {
		// 	log.Info("Message", zap.Any("message", update.Message))
		// }

		// log.Info("Message", zap.Any("message", update.Message))
		// return nil

		// Don't echo service message.
		msg, ok := update.Message.(*tg.Message)
		if !ok {
			return nil
		}

		// Echo received message.
		_, err := sender.Answer(e, update).Text(ctx, msg.Message)
		return err
	})

	tgc.client.Run(tgc.ctx, func(ctx context.Context) error {
		// Perform auth if no session is available.
		// if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		// 	return errors.Wrap(err, "auth")
		// }

		// Fetch user info.
		// user, err := client.Self(ctx)
		// if err != nil {
		// 	return errors.Wrap(err, "call self")
		// }
		err = tgc.gupMsg.Run(tgc.ctx, api, self.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("Gaps started")
			},
		})

		return err
	})

	err = tgc.gupMsg.Run(tgc.ctx, api, self.ID, updates.AuthOptions{
		OnStart: func(ctx context.Context) {
			log.Info("Gaps started")
		},
	})

	if err != nil {
		log.Info(err.Error())
	}
	// msgs, err := api.

	// https://github.com/gotd/td/blob/main/examples/updates/main.go

	return nil
}
