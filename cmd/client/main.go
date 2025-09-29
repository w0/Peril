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
		log.Fatalf("could not connect to rabbitMQ: %s", err)
	}

	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not display welcome screen: %s", err)
	}

	pauseQueue := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, pauseQueue, routing.PauseKey, pubsub.TransientQueue)
	if err != nil {
		log.Fatalf("bind to pause queue failed: %s", err)
	}

	gs := gamelogic.NewGameState(username)

	movesQueue := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, movesQueue, "army_moves.*", pubsub.TransientQueue, handlerMove(gs))
	if err != nil {
		log.Fatalf("moves subscribe failed: %s", err)
	}

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, pauseQueue, routing.PauseKey, pubsub.TransientQueue, handlerPause(gs))
	if err != nil {
		log.Fatalf("pause subscribe failed: %s", err)
	}

	pubCh, _ := conn.Channel()

	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Printf("spawn error: %s\n", err)
			}
		case "move":
			move, err := gs.CommandMove(words)
			if err != nil {
				fmt.Printf("unit move error: %s\n", err)
				continue
			}

			err = pubsub.PublishJSON(pubCh, routing.ExchangePerilTopic, movesQueue, move)
			if err != nil {
				fmt.Printf("publish move error: %s", err)
				continue
			}
			fmt.Printf("move %s %d", move.ToLocation, len(move.Units))
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Printf("unknown command: %s", words[0])
		}
	}

}
