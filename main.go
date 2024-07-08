package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/EfimoffN/beholder/config"
	"github.com/EfimoffN/beholder/kfkapi"
	"github.com/EfimoffN/beholder/tg_beholder"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Service beholder run")

	config, err := config.CreateConfig()
	if err != nil {
		log.Error().Err(err)

		return
	}

	tgClient := tg_beholder.CreateTgBeholder( // получать настрйоки аккаунта из БД
		config.BeholderTG.PhoneNumber,
		config.BeholderTG.AppHASH,
		config.BeholderTG.SessionTG,
		config.BeholderTG.AppID,
		config.BeholderTG.SessionOptMin,
		config.BeholderTG.SessionOptMax,
		config.BeholderTG.CapChan,
		ctx,
	)

	err = tgClient.Authorize()
	if err != nil {
		log.Error().Err(err)

		return
	}

	kfk, err := kfkapi.CreateKafkaProducer([]string{config.ConfigKfk.ProducerBroker}, log, config.ConfigKfk.ProducerTopic)
	if err != nil {
		log.Error().Err(err)

		return
	}

	wrk := CreateWork(log, kfk, &tgClient)

	err = wrk.WorkerFunc(config.ConfigKfk.ProducerTopic, ctx)
	if err != nil {
		log.Error().Err(err).Msg("process work failed")

		return
	}

	log.Debug().Msg("process work finished successfully")
}
