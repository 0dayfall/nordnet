package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/0dayfall/nordnet"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("Welcome to (N)ordnet (T)rading (A)lgorithm...")
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() error {
		publicFeed := nordnet.NewPublicFeed("https://api.test.nordnet.se/next/2")
		msgChan := make(chan *nordnet.PublicMsg)
		errChan := make(chan error)
		publicFeed.Dispatch(msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				logger.Info(msg)
			case err := <-errChan:
				logger.Error(err)
				return err
			case <-ctx.Done():
				logger.Info("Shutdown public feed")
				return nil
			}
		}
	}()

	go func() error {
		privateFeed := nordnet.NewPrivateFeed("https://api.test.nordnet.se/next/2")
		msgChan := make(chan *nordnet.PrivateMsg)
		errChan := make(chan error)
		privateFeed.Dispatch(msgChan, errChan)

		for {
			select {
			case msg := <-msgChan:
				log.Println(msg)
			case err := <-errChan:
				log.Println(err)
				return err
			case <-ctx.Done():
				log.Println("Shutdown private feed")
				return nil
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Bye from Nordnet trading algorithm!")
			return
		}
	}
}
