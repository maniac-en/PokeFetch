// Package client implements data fetching for the PokeAPI (pokeapi.co)
package client

import (
	"encoding/json"
	"fmt"
	"io"
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
	requestURL := mapAreaEndpoint + fmt.Sprintf("/%s", *mapAreaName)
	return GetResourceFromPokeAPI[PokeMapArea](client, &requestURL)
}

func (client *Client) GetPokemon(pokemonName *string) (Pokemon, error) {
	requestURL := pokemonEndpoint + fmt.Sprintf("/%s", *pokemonName)
	return GetResourceFromPokeAPI[Pokemon](client, &requestURL)
}

func GetResourceFromPokeAPI[T any](client *Client, URL *string) (T, error) {
	if URL == nil {
		return *new(T), fmt.Errorf("error fetching an empty resource")
	}

	var result T
	if val, ok := client.cache.Get(*URL); ok {
		fmt.Println("Providing cached result from past", client.cache.GetTTL(), "minutes")
		if err := json.Unmarshal(val, &result); err != nil {
			return *new(T), err
		}
		return result, nil
	}

	req, err := http.NewRequest(http.MethodGet, *URL, nil)
	if err != nil {
		return *new(T), nil
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return *new(T), nil
	}
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return *new(T), fmt.Errorf("resource not found")
		}
		return *new(T), fmt.Errorf("received the response with %v status", res.StatusCode)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return *new(T), err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return *new(T), err
	}
	client.cache.Add(*URL, data)
	return result, nil
}
