package main

import (
	"time"

	"github.com/maniac-en/pokefetch/internal/client"
)

func main() {
	pokeClient := client.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		client: pokeClient,
	}
	ReplStart(cfg)
}
