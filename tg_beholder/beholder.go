package tg_beholder

import (
	"context"
	"strconv"

	"github.com/EfimoffN/beholder/types"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/updates"
	"github.com/pkg/errors"

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

		ch, ok := pub.PeerID.(*tg.PeerChannel)
		if !ok {
			return nil
		}

		channel, err := tgb.SerchChannelByID(ch.ChannelID)
		if err != nil {
			return err
		}

		// добавить поиск каналов по id
		// добавлять ID родительскокго канала, ID группы комментов
		// Id публикации
		// время публикации
		// Id для комментария в чате (если есть)

		acceptedPublication := types.AcceptedPublication{
			MessageID:   int64(pub.ID),
			ChannelTgID: ch.ChannelID,
			MessageLink: "https://t.me/" + channel.Username + "/" + strconv.Itoa(pub.ID),
		}

		tgb.PostSend <- acceptedPublication

		return nil
	})

	// Create message sending helper.
	sender := message.NewSender(tgb.client.API())

	tgb.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		// Don't echo service message.
		msg, ok := update.Message.(*tg.Message)
		if !ok {
			return nil
		}

		// Echo received message.
		_, err := sender.Answer(e, update).Text(ctx, msg.Message)
		return err
	})

	err = tgb.client.Run(tgb.ctx, func(ctx context.Context) error {
		err = tgb.gupMsg.Run(tgb.ctx, api, self.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("Gaps started")
			},
		})

		return err
	})

	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = tgb.gupMsg.Run(tgb.ctx, api, self.ID, updates.AuthOptions{
		OnStart: func(ctx context.Context) {
			log.Info("Gaps started")
		},
	})

	if err != nil {
		log.Info(err.Error())
		return err

	}

	return nil
}

func (tgb *TgBeholder) SerchChannelByID(channelID int64) (*tg.Channel, error) {
	api := tgb.client.API()

	req := []tg.InputChannelClass{
		&tg.InputChannel{
			ChannelID: channelID,
		},
	}

	channelList, err := api.ChannelsGetChannels(tgb.ctx, req)
	if err != nil {
		return nil, err
	}

	if channelList.Zero() {
		err := errors.New("channel not found")

		return nil, err
	}

	channelData, ok := channelList.GetChats()[0].(*tg.Channel)
	if !ok {
		return nil, errors.New("can not convert to channel")
	}

	return channelData, nil
}
