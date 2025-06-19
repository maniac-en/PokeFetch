// Package client implements data fetching for the PokeAPI (pokeapi.co)
package client

import (
	"net/http"
	"time"

	"github.com/maniac-en/pokefetch/internal/cache"
)

type Client struct {
	cache      cache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		cache: cache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}
