package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerGameLog(gl routing.GameLog) pubsub.AckType {
	defer fmt.Println("> ")

	err := gamelogic.WriteLog(gl)
	if err != nil {
		log.Printf("error writing gamelog: %s\n", err)
		return pubsub.NackDiscard
	}

	return pubsub.Ack
}
