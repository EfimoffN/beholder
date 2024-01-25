package main

import (
	"context"
	"os"

	"github.com/EfimoffN/beholder/config"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Service run")

	config, err := config.CreateConfig()
	if err != nil {
		log.Error().Err(err)

		return
	}

	// psgConnect, err := connectDB(config.ConfigDB)
	// if err != nil {
	// 	log.Error().Err(err)

	// 	return
	// }
	// defer psgConnect.Close()

	wrk := CreateWork(log, config.SessionTG, context.Background())

	// wrk := CreateWork(apiDB, log, config.SessionTG, context.Background(), kfk)

	err = wrk.Work()
	if err != nil {
		log.Error().Err(err).Msg("process work failed")

		return
	}

	log.Debug().Msg("process work finished successfully")

}
