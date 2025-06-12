package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/0dayfall/nordnet/api"
	"github.com/0dayfall/nordnet/feed"
	"github.com/0dayfall/nordnet/util"
)

var defaultAPI = "https://api.nordnet.se/next/2"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("Welcome to (N)ordnet (T)rading (A)lgorithm...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info().Msg("Listening for SIGTERM or Interrupt...")

	if err := godotenv.Load("config/config.env"); err != nil {
		log.Warn().Msg(".env file not found, using environment variables")
	}

	apiURL := defaultAPI
	if val, ok := os.LookupEnv("NORDNET_API_URL"); ok {
		apiURL = val
		log.Info().Str("url", apiURL).Msg("Using custom API URL")
	} else {
		log.Warn().Str("default", apiURL).Msg("NORDNET_API_URL not set, using default")
	}

	username, ok := os.LookupEnv("NORDNET_USER")
	if !ok {
		log.Fatal().Msg("NORDNET_USER not set")
	}

	password, ok := os.LookupEnv("NORDNET_PASSWORD")
	if !ok {
		log.Fatal().Msg("NORDNET_PASSWORD not set")
	}

	pemPath, ok := os.LookupEnv("NORDNET_PEM_FILE")
	if !ok {
		log.Fatal().Msg("NORDNET_PEM_FILE not set")
	}

	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read PEM file")
	}

	cred, err := util.GenerateCredentials([]byte(username), []byte(password), pemData)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate credentials")
	}

	// PUBLIC FEED
	go util.SafeGo("publicFeed", func() {
		publicFeed, err := feed.NewPublicFeed(apiURL)
		if err != nil {
			log.Error().Err(err).Msg("failed to create public feed")
			stop()
			return
		}

		if err := publicFeed.Login(cred, nil); err != nil {
			log.Error().Err(err).Msg("public feed login failed")
			stop()
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		msgChan := make(chan *feed.PublicMsg)
		errChan := make(chan error)
		publicFeed.Dispatch(ctx, msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				switch msg.Type {
				case "heartbeat":
					log.Debug().Str("type", msg.Type).Msg("public heartbeat")
				case "price":
					log.Info().Interface("price", msg.Data).Msg("price update")
				case "trade":
					log.Info().Interface("trade", msg.Data).Msg("trade update")
				case "depth":
					log.Info().Interface("depth", msg.Data).Msg("depth update")
				case "trading_status":
					log.Info().Interface("status", msg.Data).Msg("trading status")
				case "indicator":
					log.Info().Interface("indicator", msg.Data).Msg("indicator update")
				case "news":
					log.Info().Interface("news", msg.Data).Msg("news update")
				}
			case err := <-errChan:
				log.Error().Err(err).Msg("public feed error")
				stop()
				return
			case <-ctx.Done():
				log.Info().Msg("shutdown: public feed")
				return
			}
		}
	})

	// PRIVATE FEED
	go util.SafeGo("privateFeed", func() {
		privateFeed, err := feed.NewPrivateFeed(apiURL)
		if err != nil {
			log.Error().Err(err).Msg("failed to create private feed")
			stop()
			return
		}

		if err := privateFeed.Login(cred, nil); err != nil {
			log.Error().Err(err).Msg("private feed login failed")
			stop()
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		msgChan := make(chan *feed.PrivateMsg)
		errChan := make(chan error)
		privateFeed.Dispatch(ctx, msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				log.Info().Interface("private", msg).Msg("private feed message")
			case err := <-errChan:
				log.Error().Err(err).Msg("private feed error")
				stop()
				return
			case <-ctx.Done():
				log.Info().Msg("shutdown: private feed")
				return
			}
		}
	})

	// API CLIENT
	go util.SafeGo("apiClient", func() {
		client := api.NewAPIClient(cred)
		client.SystemStatus()
		client.Login()

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("shutdown: API client")
				return
			}
		}
	})

	// Main shutdown listener
	<-ctx.Done()
	log.Info().Msg("Bye from NTA!")
}
