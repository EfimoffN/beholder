package tg_beholder

import (
	"context"
	"math/rand"
	"time"

	"github.com/EfimoffN/beholder/types"
	"github.com/gotd/td/telegram/updates"
	"github.com/pkg/errors"

	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

var hardLimit int = 100
var errLeadToForm = errors.New("could not lead to the form")

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
		log.Debug("Channel message", zap.Any("message", update.Message))

		pub, ok := update.Message.(*tg.Message)
		if !ok {
			return nil
		}

		_, ok = pub.FromID.(*tg.PeerUser)
		if ok {
			return nil
		}

		ch, ok := pub.PeerID.(*tg.PeerChannel)
		if !ok {
			return nil
		}

		var chat = tg.Channel{}
		for _, ent := range e.Channels {
			if ent.ID == ch.ChannelID {
				err := tgb.markMessageRead(pub.ID, ch.ChannelID, ent.AsInputPeer())
				if err != nil {
					log.Error(err.Error())

					return nil
				}

				chat = *ent
			}
		}

		if chat.ID == 0 && chat.AccessHash == 0 {
			return nil
		}

		acceptedPublications, err := tgb.GetLastPublication(&chat, hardLimit)
		if err != nil {
			return nil
		}

		for _, pub := range *acceptedPublications {
			tgb.PostSend <- pub
		}

		return nil
	})

	err = tgb.client.Run(tgb.ctx, func(ctx context.Context) error {
		err = tgb.gupMsg.Run(tgb.ctx, api, self.ID, updates.AuthOptions{
			OnStart: func(ctx context.Context) {
				log.Info("tgb.client Gaps started")
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
			log.Info("tgb.gupMsg Gaps started")
		},
	})

	if err != nil {
		log.Info(err.Error())
		return err

	}

	return nil
}

func (tgb *TgBeholder) GetLastPublication(channelData *tg.Channel, limit int) (*[]types.AcceptedPublication, error) {
	result := make([]types.AcceptedPublication, 0, limit)

	publications, err := tgb.GetChannelPublication(channelData, limit)
	if err != nil {
		return nil, err
	}

	for _, pub := range publications.Messages {
		messagePeer, ok := pub.(*tg.Message)
		if ok {
			acceptedPublication := types.AcceptedPublication{
				ChannelTgID:      channelData.ID,
				MessageChannelID: int64(messagePeer.ID),
				CreatedDate:      int64(messagePeer.Date),
				EditDate:         int64(messagePeer.EditDate),
			}

			result = append(result, acceptedPublication)
		}

	}

	return &result, nil
}

func (tgb *TgBeholder) GetChannelPublication(
	channelData *tg.Channel,
	limit int,
) (*tg.MessagesChannelMessages, error) {
	// новые аккаунты разработчиков имеют ограничение не более 30 обращений к серверу в минуту
	// ставим задержку что бы не получить бан
	// time.Sleep(time.Duration(RandInt(tgb.sessionOptMin, tgb.sessionOptMax)) * time.Millisecond)
	api := tgb.client.API()

	messagesRequest := &tg.MessagesGetHistoryRequest{
		Limit: hardLimit,
		Peer:  channelData.AsInputPeer(),
		Hash:  channelData.AsInput().AccessHash,
	}

	if limit != 0 {
		messagesRequest.Limit = limit
	}

	res, err := api.MessagesGetHistory(tgb.ctx, messagesRequest)
	if err != nil {
		return nil, errors.Wrap(err, "messages get history failed")
	}

	channelPosts, ok := res.(*tg.MessagesChannelMessages)
	if !ok {
		return nil, errors.Wrap(errLeadToForm, "could not lead to the form messages channel messages")
	}

	return channelPosts, nil
}

func (tgb *TgBeholder) markMessageRead(messageId int, channelID int64, peer tg.InputPeerClass) error {
	api := tgb.client.API()

	mrkReq := tg.ChannelsReadHistoryRequest{
		MaxID: messageId,
		Channel: &tg.InputChannelFromMessage{
			ChannelID: channelID,
			MsgID:     messageId,
			Peer:      peer,
		},
	}

	mrk, err := api.ChannelsReadHistory(tgb.ctx, &mrkReq)
	if !mrk {
		return nil
	}

	if err != nil {
		tgb.Logger.Error().Err(err)

		return err
	}

	time.Sleep(time.Duration(RandInt(2000, 4000)) * time.Millisecond)
	return nil
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
