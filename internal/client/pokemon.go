package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetPokemon(pokemonName *string) (Pokemon, error) {
	if pokemonName == nil {
		return Pokemon{}, fmt.Errorf("empty pokemon name allowed")
	}
	requestURL := pokemonEndpoint + fmt.Sprintf("/%s", *pokemonName)

	if val, ok := c.cache.Get(requestURL); ok {
		fmt.Println("Providing cached result from past", c.cache.GetTTL(), "minutes")
		pokemon := Pokemon{}
		if err := json.Unmarshal(val, &pokemon); err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return Pokemon{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}

	if res.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("received the response with %v status", res.StatusCode)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	pokemon := Pokemon{}
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return Pokemon{}, err
	}
	c.cache.Add(requestURL, data)
	return pokemon, nil
}
