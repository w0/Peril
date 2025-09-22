package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	rUrl := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rUrl)
	if err != nil {
		fmt.Printf("failed to connect %w", err)
		os.Exit(1)
	}

	defer conn.Close()

	login, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("failed creating login %w", err)
		os.Exit(1)
	}

	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, fmt.Sprintf("%s.%s", routing.PauseKey, login), routing.PauseKey, pubsub.TransientQueue)
	if err != nil {
		fmt.Printf("failed binding queue %s", err)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("failed binding queue %s", err)
		os.Exit(1)
	}

	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// Add this to keep the program running so you can see the queue:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")

}
