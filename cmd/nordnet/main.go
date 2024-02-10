package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/0dayfall/nordnet/feed"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Welcome to (N)ordnet (T)rading (A)lgorithm...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() error {
		publicFeed, err := feed.NewPublicFeed("https://api.test.nordnet.se/next/2")
		if err != nil {
			logger.Error("Error: %v", err)
		}
		msgChan := make(chan *feed.PublicMsg)
		errChan := make(chan error)
		publicFeed.Dispatch(msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				logger.Info("private feed: %v", msg)
			case err := <-errChan:
				logger.Error("Error: %v", err)
				return err
			case <-ctx.Done():
				logger.Info("Shutdown public feed")
				return nil
			}
		}
	}()

	go func() error {
		privateFeed, err := feed.NewPrivateFeed("https://api.test.nordnet.se/next/2")
		if err != nil {
			logger.Error("Error: %v", err)
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

	for {
		select {
		case <-ctx.Done():
			logger.Info("Bye from Nordnet trading algorithm!")
			return
		}
	}
}
