package main

import (
	"context"
	"errors"
	"os"

	beholder "github.com/EfimoffN/beholder/tg_beholder"
	"github.com/EfimoffN/beholder/types"
	"github.com/rs/zerolog"
)

var errAlreadyExists = errors.New("can't create a file that already exists")
var errMsgConvert = errors.New("can't convert message")

const (
	sessionOptMin = 2000
	sessionOptMax = 3000
)

type Worker struct {
	// db   *sqlapi.API
	log  zerolog.Logger
	UUID string
	ctx  context.Context
	// kfk  *kfkapi.KafkaProducer
	// kafka & nsq
}

func CreateWork(
	// db *sqlapi.API,
	log zerolog.Logger,
	uuid string,
	ctx context.Context,
	// kfk *kfkapi.KafkaProducer,
) *Worker {
	wrk := Worker{
		// db:   db,
		log:  log,
		UUID: uuid,
		ctx:  ctx,
		// kfk:  kfk,
	}

	return &wrk
}

func (w *Worker) Work() error {
	// configTG, err := w.db.GetSessionByID(w.UUID)
	// if err != nil {
	// 	w.log.Error().Err(err)
	// 	return err
	// }

	configTG := types.SessionRow{}

	// lastMsgs, err := w.db.GetLastMSGs(w.UUID)
	// if err != nil {
	// 	w.log.Error().Err(err)
	// 	return err
	// }

	// if lastMsgs == nil || len(*lastMsgs) == 0 {
	// 	return nil // ошибка которую обработать?
	// }

	fileStorage, err := w.createSession(configTG.Session)
	if err != nil {
		w.log.Error().Err(err)
		return err
	}

	tgCrawler := beholder.InitTgCrawler(w.log, configTG.PhoneNumber, fileStorage, configTG.AppID, configTG.AppHash, w.ctx, sessionOptMin, sessionOptMax)

	err = tgCrawler.Authorize()
	if err != nil {
		w.log.Error().Err(err)
		return err
	}

	err = tgCrawler.CheckedPosts()
	if err != nil {
		w.log.Error().Err(err)
	}

	// go func() {
	// 	w.kfk.Sender()
	// }()

	// for _, msg := range *lastMsgs {
	// 	channelProp, err := w.db.GetChannelByTgID(msg.ChanneTgID)
	// 	if err != nil {
	// 		w.log.Error().Err(err)
	// 		continue
	// 	}

	// 	channel, err := tgCrawler.SearchChannel(channelProp.ChannelName)
	// 	if err != nil {
	// 		w.log.Error().Err(err)
	// 		continue
	// 	}

	// 	currentMsgs := &tg.MessagesChannelMessages{
	// 		Flags:          0,
	// 		Inexact:        false,
	// 		Pts:            0,
	// 		Count:          0,
	// 		OffsetIDOffset: 0,
	// 		Messages:       []tg.MessageClass{},
	// 		Topics:         []tg.ForumTopicClass{},
	// 		Chats:          []tg.ChatClass{},
	// 		Users:          []tg.UserClass{},
	// 	}

	// 	if msg.LastMsgID == 0 {
	// 		currentMsgs, err = tgCrawler.GetChannelPublications(channel, 0, 0, 0, 0, 1)
	// 		if err != nil {
	// 			w.log.Error().Err(err)
	// 			continue
	// 		}
	// 	} else {
	// 		currentMsgs, err = beholder.GetChannelPublications(channel, 0, 0, 0, int(msg.LastMsgID), 100)
	// 		if err != nil {
	// 			w.log.Error().Err(err)
	// 			continue
	// 		}

	// 		err = w.addQueue(*currentMsgs, channelProp.ChannelName, channelProp.ChannelTgID)
	// 		if err != nil {
	// 			w.log.Error().Err(err)
	// 			continue
	// 		}

	// 	}

	// 	if len(currentMsgs.Messages) == 0 {
	// 		continue
	// 	}

	// 	msg, ok := currentMsgs.Messages[0].(*tg.Message)
	// 	if !ok {
	// 		w.log.Error().Err(errMsgConvert)
	// 		continue
	// 	}

	// 	err = w.saveLastMsg(msg, channelProp.ChannelTgID) // сообщения в порядке убывания
	// 	if err != nil {
	// 		w.log.Error().Err(err)
	// 	}
	// }

	// defer w.kfk.ProducerClose()
	// defer w.kfk.WG.Wait()

	return nil
}

func (w *Worker) createSession(session string) (string, error) {
	// fileStorage := strconv.FormatInt(time.Now().UTC().Unix(), 10) + ".json"

	// sessionTG := &beholder.SessionTG{}

	// err := json.Unmarshal([]byte(session), sessionTG)
	// if err != nil {
	// 	w.log.Error().Err(err).Msg("unmarshal session data")

	// 	return "", err
	// }

	// sessionData, err := json.Marshal(sessionTG)
	// if err != nil {
	// 	w.log.Debug().Err(err).Msg("marshal session property")

	// 	return "", err
	// }

	// err = w.createFileSession(fileStorage, sessionData)
	// if err != nil {
	// 	w.log.Error().Err(err).Msg("create session err")

	// 	return "", err
	// }
	pathToSessionFile := ""
	fileStorage := pathToSessionFile

	return fileStorage, nil
}

func (w *Worker) createFileSession(name string, val []byte) error {
	_, err := os.Stat(name)
	if os.IsExist(err) {
		w.log.Error().Err(errAlreadyExists).Str("file", name)

		return errAlreadyExists
	}

	file, err := os.Create(name)
	if err != nil {
		w.log.Error().Err(err).Str("file", name).Msg("failed create file session")

		return err
	}
	defer file.Close()

	_, err = file.Write(val)
	if err != nil {
		w.log.Error().Err(err).Str("file", name).Msg("failed write in to file session")

		return err
	}

	return nil
}

// func (w *Worker) saveLastMsg(msg *tg.Message, idChn int64) error {
// 	err := w.db.UpdateLastMSG(idChn, msg.GetID())
// 	if err != nil {
// 		w.log.Error().Err(err)
// 		return err
// 	}

// 	return nil
// }

// func (w *Worker) addQueue(msgs tg.MessagesChannelMessages, channelName string, channelTgID int64) error {
// 	var groupMsg int64
// 	var headMsg *tg.Message

// 	for i := len(msgs.Messages) - 1; i >= 0; i-- {
// 		msgCls := msgs.Messages[i]
// 		msg, ok := msgCls.(*tg.Message)
// 		if !ok {
// 			w.log.Error().Err(errMsgConvert)
// 			continue
// 		}

// 		if msg.GroupedID != 0 {
// 			headMsg = msg
// 			groupMsg = msg.GroupedID

// 			continue
// 		}

// 		if groupMsg != 0 && msg.GroupedID == 0 {
// 			groupMsg = 0

// 			obj := types.ReadyRepost{
// 				ChannelTgID: channelTgID,
// 				MessageLink: "https://t.me/" + channelName + "/" + strconv.Itoa(headMsg.ID),
// 			}

// 			w.kfk.WG.Add(1)
// 			w.kfk.MessegSend <- obj
// 		}

// 		obj := types.ReadyRepost{
// 			ChannelTgID: channelTgID,
// 			MessageLink: "https://t.me/" + channelName + "/" + strconv.Itoa(msg.ID),
// 		}

// 		w.kfk.WG.Add(1)
// 		w.kfk.MessegSend <- obj
// 	}

// 	return nil
// }
