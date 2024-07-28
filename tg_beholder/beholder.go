package tg_beholder

import (
	"context"
	"strconv"

	"github.com/EfimoffN/beholder/types"
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

		// получаем чат в который пересланно собщениеc
		// ищем этот чат и ищем сообщение по Id сообщения и по Id группы в которой опулдикованно сообщениеctx
		// получаем id обсуждения и отправляем дальше для комментраиев

		ch, ok := pub.PeerID.(*tg.PeerChannel)
		if !ok {
			return nil
		}

		if pub.Replies.ChannelID == 0 {
			return nil
		}

		accessHash, err := tgb.SerchChannelByID(ch.ChannelID, pub.Replies.ChannelID, pub.Message)
		if err != nil {
			return err
		}

		acceptedPublication := types.AcceptedPublication2{
			ChannelTgID:      ch.ChannelID,
			ChatTgID:         pub.Replies.ChannelID,
			MessageChannelID: int64(pub.ID),
			MessageChatID:    int64(accessHash), проверить получение верности ID чата
			Created:          int64(pub.Date),
			TextMessage:      pub.Message,
		}

		tgb.PostSend <- acceptedPublication

		return nil
	})

	// Create message sending helper.
	// sender := message.NewSender(tgb.client.API())

	tgb.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		// Don't echo service message.
		msg, ok := update.Message.(*tg.Message)
		if !ok {
			return nil
		}

		if msg != nil {
			return nil
		}

		// Echo received message.
		// _, err := sender.Answer(e, update).Text(ctx, msg.Message)
		return nil
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

func (tgb *TgBeholder) SerchChannelByID(channelID int64, peerChannelId int64, textMsg string) (int64, error) {
	api := tgb.client.API()
	var accessHash int64 = 0

	req := []tg.InputChannelClass{
		&tg.InputChannel{
			ChannelID: channelID,
		},
	}

	channelList, err := api.ChannelsGetChannels(tgb.ctx, req)
	if err != nil {
		return accessHash, err
	}

	if channelList.Zero() {
		err := errors.New("channel not found")

		return accessHash, err
	}

	channelData, ok := channelList.GetChats()[0].(*tg.Channel)
	if !ok {
		return accessHash, errors.New("can not convert to channel")
	}

	chFull := tg.InputChannel{
		ChannelID:  channelData.ID,
		AccessHash: channelData.AsInput().AccessHash,
	}

	resFull, err := api.ChannelsGetFullChannel(tgb.ctx, &chFull)
	if err != nil {
		err = errors.Wrapf(err, "can't find the channel by ID: '%s'", strconv.FormatInt(channelID, 10))
		return accessHash, err
	}

	for _, chat := range resFull.GetChats() {
		if chat.GetID() == peerChannelId {
			chA, ok := chat.(*tg.Channel)
			if ok {
				accessHash = chA.AccessHash

				messagesRequest := &tg.MessagesGetHistoryRequest{
					Limit: 100,
					Peer:  chA.AsInputPeer(),
					Hash:  accessHash,
					// MaxID: msgID + 1,
					// MinID: msgID - 1,
				}

				messages, err := api.MessagesGetHistory(tgb.ctx, messagesRequest)
				if err != nil {
					return 0, errors.New("can not get messages")
				}

				if messages == nil {
					return accessHash, nil
				}

				msg, ok := messages.(*tg.MessagesChannelMessages)
				if ok {
					for _, message := range msg.Messages {
						messagePeer, ok := message.(*tg.Message)
						if ok {
							messageF := messagePeer.FromID
							from, ok := messageF.(*tg.PeerChannel)
							if ok {
								if from.ChannelID == channelID {
									messageP := messagePeer.PeerID
									peer, ok := messageP.(*tg.PeerChannel)
									if ok {
										return peer.ChannelID, nil
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if accessHash == 0 {
		return accessHash, errors.New("can not found comments chat")
	}

	return accessHash, nil
}

func (tgb *TgBeholder) getRepliesMessageChat(accessHash int64, chatId int64) (int, error) {
	api := tgb.client.API()

	peer := tg.InputPeerChat{
		ChatID: chatId,
	}

	messagesRequest := &tg.MessagesGetHistoryRequest{
		Limit: 100,
		Peer:  &peer,
		Hash:  accessHash,
		// MaxID: msgID + 1,
		// MinID: msgID - 1,
	}

	messages, err := api.MessagesGetHistory(tgb.ctx, messagesRequest)
	if err != nil {
		return 0, errors.New("can not get messages")
	}

	channelPosts, ok := messages.(*tg.MessagesChannelMessages)
	if !ok {
		return 0, errors.New("could not lead to the form messages channel messages")
	}

	for _, message := range channelPosts.Messages {
		if message.GetID() == 0 {
			return message.GetID(), nil
		}
	}

	return 0, errors.New("not found message")
}
