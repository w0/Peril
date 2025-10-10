package gamelogic

import (
	"errors"
	"strconv"
)

func (gs *GameState) CommandSpam(words []string) (int, error) {
	if len(words) < 2 {
		return 0, errors.New("usage: spam <number>")
	}

	numSpams, err := strconv.Atoi(words[1])
	if err != nil {
		return 0, err
	}

	return numSpams, nil

}
