package main

import (
	"fmt"
	"time"

	"github.com/maniac-en/pokefetch/internal/client"
)

func main() {
	pokeClient, err := client.NewClient(5*time.Second, 1*time.Minute)
	if err != nil {
		panic(fmt.Sprintf("error creating a client: %v", err))
	}
	cfg := &config{
		client:  *pokeClient,
		pokedex: make(map[string]client.Pokemon),
	}
	ReplStart(cfg)
}
