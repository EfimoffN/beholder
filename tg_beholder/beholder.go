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

		_, ok = pub.FromID.(*tg.PeerUser)
		if ok {
			return nil
		}

		ch, ok := pub.PeerID.(*tg.PeerChannel)
		if !ok {
			return nil
		}

		if pub.Replies.ChannelID == 0 {
			return nil
		}

		messageIDChat, err := tgb.SerchChannelByID(ch.ChannelID, pub.Replies.ChannelID, pub.Message)
		if err != nil {
			return err
		}

		acceptedPublication := types.AcceptedPublication2{
			ChannelTgID:      ch.ChannelID,
			ChatTgID:         pub.Replies.ChannelID,
			MessageChannelID: int64(pub.ID),
			MessageChatID:    messageIDChat,
			Created:          int64(pub.Date),
			TextMessage:      pub.Message,
		}

		tgb.PostSend <- acceptedPublication

		return nil
	})

	// tgb.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
	// 	// Don't echo service message.
	// 	msg, ok := update.Message.(*tg.Message)
	// 	if !ok {
	// 		return nil
	// 	}

	// 	if msg != nil {
	// 		return nil
	// 	}

	// 	return nil
	// })

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

	chFull := tg.InputChannel{}
	for _, channelData := range channelList.GetChats() {
		channelD, ok := channelData.(*tg.Channel)
		if ok && channelD.GetID() == channelID {
			chFull.ChannelID = channelD.ID
			chFull.AccessHash = channelD.AsInput().AccessHash
		}
	}
	if chFull.AccessHash == 0 {
		return accessHash, errors.New("can not convert to channel")
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
									return int64(messagePeer.ID), nil
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

// func (tgb *TgBeholder) markMessageRead(messageId int64, peerChannel ) error{

// }
