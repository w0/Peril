package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	rUrl := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rUrl)
	if err != nil {
		fmt.Printf("failed opening connection %w", err)
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Printf("amqp connection succesfull.\n")

	pubChan, err := conn.Channel()
	if err != nil {
		fmt.Printf("failed creating channel %w", err)
		os.Exit(1)
	}

	err = pubsub.PublishJSON(pubChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})

	if err != nil {
		fmt.Printf("failed to publish json %w", err)
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Printf("interrupt recived. Shutting down server..\n")

}
