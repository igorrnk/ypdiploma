package main

import (
	"context"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Log().Msg("Gophermart is running")

	config, err := configs.InitConfig()
	if err != nil {
		log.Fatal().Msgf("Config error: %v", err)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFieldFormat
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Info().Msgf("Config: %+v: ", config)

	gophermart, err := service.NewGophermart(context.Background(), config)
	if err != nil {
		log.Fatal().Msgf("InitGophermart error: %v", err)
	}
	if err = gophermart.Run(); err != nil {
		log.Fatal().Msgf("Gophermart stopped with error: %v", err)
	}

}
