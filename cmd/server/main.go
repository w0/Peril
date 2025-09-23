package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to rabbitMQ: %w", err)
	}

	defer conn.Close()

	fmt.Printf("peril game server created.\n")

	pubChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %w", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Printf("sending pause message..\n")
			err = pubsub.PublishJSON(pubChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})

			if err != nil {
				log.Fatalf("could not publish: %w", err)
			}
		case "resume":
			fmt.Printf("sending resume message..\n")
			err = pubsub.PublishJSON(pubChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})

			if err != nil {
				log.Fatalf("could not publish: %w", err)
			}
		case "quit":
			log.Println("exiting server.")
			return
		default:
			fmt.Printf("unknown command %s", words[0])

		}
	}
}
