// Package client implements data fetching for the PokeAPI (pokeapi.co)
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/maniac-en/pokefetch/internal/cache"
)

type Client struct {
	cache      cache.Cache
	httpClient http.Client
}

func NewClient(timeout, cacheInterval time.Duration) (*Client, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	if cacheInterval <= 0 {
		return nil, fmt.Errorf("cache interval must be positive")
	}
	return &Client{
		cache: cache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}, nil
}

func (client *Client) GetMapAreas(pageURL *string) (PokeMapAreas, error) {
	// if the passed URL is empty, i.e., the next/prev URL passed down is empty,
	// then use the default endpoint which fetches the first page
	if pageURL != nil {
		return GetResourceFromPokeAPI[PokeMapAreas](client, pageURL)
	}
	defaultURL := mapAreaDefaultEndpoint
	return GetResourceFromPokeAPI[PokeMapAreas](client, &defaultURL)
}

func (client *Client) GetMapArea(mapAreaName *string) (PokeMapArea, error) {
	baseURL, _ := url.Parse(mapAreaEndpoint)
	requestURL := baseURL.JoinPath(*mapAreaName).String()
	return GetResourceFromPokeAPI[PokeMapArea](client, &requestURL)
}

func (client *Client) GetPokemon(pokemonName *string) (Pokemon, error) {
	baseURL, _ := url.Parse(pokemonEndpoint)
	requestURL := baseURL.JoinPath(*pokemonName).String()
	return GetResourceFromPokeAPI[Pokemon](client, &requestURL)
}

func GetResourceFromPokeAPI[T any](client *Client, URL *string) (T, error) {
	var zero T
	if URL == nil {
		return zero, fmt.Errorf("request URL cannot be empty")
	}

	var result T
	if val, ok := client.cache.Get(*URL); ok {
		fmt.Println("Providing cached result from past", client.cache.GetTTL(), "minutes")
		if err := json.Unmarshal(val, &result); err != nil {
			return zero, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}
		return result, nil
	}

	req, err := http.NewRequest(http.MethodGet, *URL, nil)
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("failed to execute request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return zero, fmt.Errorf("resource not found at %s", *URL)
		}
		return zero, fmt.Errorf("received the response with %v status", res.StatusCode)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	client.cache.Add(*URL, data)
	return result, nil
}
