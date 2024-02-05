package main

import (
	"context"
	"errors"

	"github.com/EfimoffN/beholder/kfkapi"
	"github.com/EfimoffN/beholder/tg_beholder"
	"github.com/rs/zerolog"
)

var errAlreadyExists = errors.New("can't create a file that already exists")
var errMsgConvert = errors.New("can't convert message")

const (
	sessionOptMin = 2000
	sessionOptMax = 3000
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

}
