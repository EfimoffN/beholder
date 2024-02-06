package main

import (
	"context"

	"github.com/EfimoffN/beholder/kfkapi"
	"github.com/EfimoffN/beholder/tg_beholder"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
)

type Worker struct {
	log      zerolog.Logger
	producer *kfkapi.KafkaProducer
	beholder *tg_beholder.TgBeholder
}

func CreateWork(
	log zerolog.Logger,
	producer *kfkapi.KafkaProducer,
	beholder *tg_beholder.TgBeholder,
) *Worker {
	wrk := Worker{
		log:      log,
		producer: producer,
		beholder: beholder,
	}

	return &wrk
}

func (w *Worker) WorkerFunc(topicProducer string, ctx context.Context) error {
	go func() {
		w.beholder.CheckedPosts()
	}()

	for {
		select {
		case <-ctx.Done():
			w.beholder.Stop()
			w.log.Debug().Msg("stoped beholder")
			w.producer.ProducerClose()

			for msg := range w.beholder.PostSend {
				result, err := json.Marshal(msg)
				if err != nil {
					w.log.Error().Err(err)
					continue
				}

				err = w.producer.SendPost(string(result), topicProducer)
				if err != nil {
					w.log.Error().Err(err)
					continue
				}
			}

			return nil
		case msg := <-w.beholder.PostSend:
			result, err := json.Marshal(msg)
			if err != nil {
				w.log.Error().Err(err)
				continue
			}

			err = w.producer.SendPost(string(result), topicProducer)
			if err != nil {
				w.log.Error().Err(err)
			}
		}
	}
}
