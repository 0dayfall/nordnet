package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/0dayfall/nordnet/api"
	"github.com/0dayfall/nordnet/feed"
	"github.com/0dayfall/nordnet/util"
)

const (
	API = "https://api.test.nordnet.se/next/2"
)

func main() {
	fmt.Println("Welcome to (N)ordnet (T)rading (A)lgorithm...")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	logger.Info("listening for SIGTERM..")
	defer stop()

	username, ok := os.LookupEnv("NORDNET_USER")
	if !ok {
		logger.Error("NORDNET_USER not set")
		os.Exit(1)
	}

	password, ok := os.LookupEnv("NORDNET_PASSWORD")
	if !ok {
		logger.Error("NORDNET_PASSWORD not set")
		os.Exit(2)
	}

	pemPath, ok := os.LookupEnv("NORDNET_PEM")
	if !ok {
		logger.Error("NORDNET_PEM not set")
		os.Exit(3)
	}

	// Read the pem file in a byte array
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		logger.Error("Error: %v", err)
		os.Exit(4)
	}

	logger.Info("Creating credentials..")
	cred, err := util.GenerateCredentials([]byte(username), []byte(password), pemData)
	if err != nil {
		logger.Error("Error: %v", err)
		os.Exit(5)

	}

	go func() error {
		publicFeed, err := feed.NewPublicFeed(API)
		if err != nil {
			logger.Error("Error: %v", err)
		}
		loginErr := publicFeed.Login(cred, nil)
		if err != loginErr {
			logger.Error("Error: %v", loginErr)
			stop()
		}

		msgChan := make(chan *feed.PublicMsg)
		errChan := make(chan error)
		publicFeed.Dispatch(msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				switch msg.Type {
				case "heartbeat":
					logger.Info("public feed: %v", msg.Data.(struct{}))
				case "price":
					logger.Info("public feed: %v", msg.Data.(feed.PublicPrice))
				case "trade":
					logger.Info("public feed: %v", msg.Data.(feed.PublicTrade))
				case "depth":
					logger.Info("public feed: %v", msg.Data.(feed.PublicDepth))
				case "trading_status":
					logger.Info("public feed: %v", msg.Data.(feed.PublicTradingStatus))
				case "indicator":
					logger.Info("public feed: %v", msg.Data.(feed.PublicIndicator))
				case "news":
					logger.Info("public feed: %v", msg.Data.(feed.PublicNews))
				}
			case err := <-errChan:
				logger.Error("error: %v", err)
				return err
			case <-ctx.Done():
				logger.Info("Shutdown public feed")
				return nil
			}
		}
	}()

	go func() error {
		privateFeed, err := feed.NewPrivateFeed(API)
		if err != nil {
			logger.Error("Error: %v", err)
		}
		loginErr := privateFeed.Login(cred, nil)
		if err != loginErr {
			logger.Error("Error: %v", loginErr)
			stop()
		}

		msgChan := make(chan *feed.PrivateMsg)
		errChan := make(chan error)
		privateFeed.Dispatch(msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				logger.Info("private feed: %v", msg)
			case err := <-errChan:
				logger.Error("Error: %v", err)
				return err
			case <-ctx.Done():
				logger.Info("Shutdown private feed")
				return nil
			}
		}
	}()

	// Initialize the algorithm, the input is price changes
	// and the output is orders to be placed as comamands to the API
	// The algorithm is a state machine that can be stopped and started
	// and can be configured with different strategies
	var channel chan string

	go func(channel chan string) {
		client := api.NewAPIClient(cred)
		client.SystemStatus()
		client.Login()

		for {
			select {
			case <-ctx.Done():
				logger.Info("Shutdown API client")
				return
			}
		}
	}(channel)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Bye from NTA!")
			return
		}
	}
}
