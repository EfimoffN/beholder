package tg_beholder

import (
	"context"
	"strconv"

	"github.com/EfimoffN/beholder/types"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/updates"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

func (tgb *TgBeholder) CheckedPosts() error {
	log := zap.NewExample()

	api := tgb.client.API()

	self, err := tgb.client.Self(tgb.ctx)
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
	tgb.dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		if update.Message.String() == "" {
			log.Info("Message", zap.Any("message", update.Message))
		}

		log.Info("Channel message", zap.Any("message", update.Message))

		message := update.Message

		pub, ok := message.(*tg.Message)
		if !ok {
			return nil
		}
		// добавить поиск каналов по id
		acceptedPublication := types.AcceptedPublication{
			// ChannelTgID: update.Message,
			ChannelTgID: int64(message.GetID()),
			MessageLink: "https://t.me/" + pub.PostAuthor + "/" + strconv.Itoa(pub.ID),
		}

		tgb.PostSend <- acceptedPublication

		return nil
	})

	// Create message sending helper.
	sender := message.NewSender(tgb.client.API())

	tgb.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
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

	tgb.client.Run(tgb.ctx, func(ctx context.Context) error {
		// Perform auth if no session is available.
		// if err := client.Auth().IfNecessary(ctx, flow); err != nil {
		// 	return errors.Wrap(err, "auth")
		// }

		// Fetch user info.
		// user, err := client.Self(ctx)
		// if err != nil {
		// 	return errors.Wrap(err, "call self")
		// }
		err = tgb.gupMsg.Run(tgb.ctx, api, self.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("Gaps started")
			},
		})

		return err
	})

	err = tgb.gupMsg.Run(tgb.ctx, api, self.ID, updates.AuthOptions{
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
