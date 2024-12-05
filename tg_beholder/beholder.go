package tg_beholder

import (
	"context"
	"math/rand"
	"strconv"
	"time"

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

		if pub.Replies.ChannelID == 0 {
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
			}

			if ent.ID == pub.Replies.ChannelID {
				chat = *ent
			}
		}

		if chat.ID == 0 && chat.AccessHash == 0 {
			return nil
		}

		messageIDChat, err := tgb.SerchChatMsgID(ch.ChannelID, &chat)
		if err != nil {
			messageIDChannel, err := tgb.SerchChannelByID(ch.ChannelID, chat.ID)
			if err != nil {
				log.Debug(err.Error())
				return nil
			}
			messageIDChat = messageIDChannel
		}

		acceptedPublication := types.AcceptedPublication2{
			ChannelTgID:      ch.ChannelID,
			ChatTgID:         pub.Replies.ChannelID,
			MessageChannelID: int64(pub.ID),
			MessageChatID:    messageIDChat,
			Created:          int64(pub.Date),
			TextMessage:      pub.Message,
		}

		log.Debug("accepted publication", zap.Any("publication", acceptedPublication))

		tgb.PostSend <- acceptedPublication

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

func (tgb *TgBeholder) SerchChatMsgID(channelID int64, chat *tg.Channel) (int64, error) {
	api := tgb.client.API()

	messagesRequest := &tg.MessagesGetHistoryRequest{
		Limit: 100,
		Peer:  chat.AsInputPeer(),
		Hash:  chat.AccessHash,
	}

	messages, err := api.MessagesGetHistory(tgb.ctx, messagesRequest)
	if err != nil {
		err := errors.New("can not get messages")
		tgb.Logger.Error().Err(err)

		return 0, err
	}

	if messages == nil {
		err := errors.New("messages is nill")
		tgb.Logger.Error().Err(err)

		return 0, err
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

	err = errors.New("message not found")
	tgb.Logger.Error().Err(err)

	return 0, err
}

// // переписать поиск пбликации используя прищедшие наьоры групп и чатов с сообщением
func (tgb *TgBeholder) SerchChannelByID(channelID int64, peerChannelId int64) (int64, error) {
	api := tgb.client.API()
	var accessHash int64 = 0

	req := []tg.InputChannelClass{
		&tg.InputChannel{
			ChannelID: channelID,
		},
	}

	channelList, err := api.ChannelsGetChannels(tgb.ctx, req)
	if err != nil {
		tgb.Logger.Error().Err(err)

		return accessHash, err
	}

	if channelList.Zero() {
		err := errors.New("channel not found")
		tgb.Logger.Error().Err(err)

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
		err := errors.New("can not convert to channel")
		tgb.Logger.Error().Err(err)

		return accessHash, err
	}

	resFull, err := api.ChannelsGetFullChannel(tgb.ctx, &chFull)
	if err != nil {
		err = errors.Wrapf(err, "can't find the channel by ID: '%s'", strconv.FormatInt(channelID, 10))
		tgb.Logger.Error().Err(err)

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
					err := errors.New("can not get messages")
					tgb.Logger.Error().Err(err)

					return 0, err
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
		err := errors.New("can not found comments chat")
		tgb.Logger.Error().Err(err)

		return accessHash, err
	}

	return accessHash, nil
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
