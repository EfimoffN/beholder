package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	log.Info().Msg("Service run")

}
