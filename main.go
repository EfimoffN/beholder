package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/EfimoffN/beholder/config"
	"github.com/EfimoffN/beholder/kfkapi"
	"github.com/EfimoffN/beholder/sqlapi"
	"github.com/EfimoffN/beholder/tg_beholder"
	"github.com/EfimoffN/beholder/types"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	db, err := connectDB(config.ConfigDB)
	if err != nil {
		log.Error().Err(err)

		return
	}
	defer db.Close()

	apiStorage := sqlapi.NewAPI(db)

	sessionRow, err := apiStorage.GetSessionsByID(config.SessionTGID)
	if err != nil {
		log.Error().Err(err)

		return
	}

	tgClient, err := tg_beholder.CreateTgBeholder(
		log,
		sessionRow.PhoneNumber,
		sessionRow.AppHash,
		sessionRow.Sessiontxt,
		sessionRow.AppID,
		config.BeholderTG.SessionOptMin,
		config.BeholderTG.SessionOptMax,
		config.BeholderTG.CapChan,
		ctx,
	)
	if err != nil {
		log.Error().Err(err)

		return
	}

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

	wrk := CreateWork(log, kfk, tgClient)

	err = wrk.WorkerFunc(config.ConfigKfk.ProducerTopic, ctx)
	if err != nil {
		log.Error().Err(err).Msg("process work failed")

		return
	}

	log.Debug().Msg("process work finished successfully")
}

// connectDB ...
func connectDB(cfg types.ConfigPsg) (*sqlx.DB, error) {
	wp := "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"
	connectionString := fmt.Sprintf(
		wp,
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBname,
		cfg.SSLmode,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Error().Err(err).Msg("sqlx.Open failed")

		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Error().Err(err).Msg("DB.Ping failed")

		return nil, err
	}

	return db, nil
}
