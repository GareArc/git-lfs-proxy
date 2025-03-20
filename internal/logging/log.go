package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(level string) {
	// unix timestamp
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// pretty print
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// set log level
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(lvl)
}
