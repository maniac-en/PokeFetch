package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetMapAreas(pageURL *string) (PokeMapAreas, error) {
	var requestURL string
	if pageURL != nil {
		requestURL = *pageURL
	} else {
		requestURL = mapAreaEndpoint + "?offset=0&limit=20"
	}

	if val, ok := c.cache.Get(requestURL); ok {
		fmt.Println("Providing cached result from past", c.cache.GetTTL(), "minutes")
		pokeMapAreas := PokeMapAreas{}
		if err := json.Unmarshal(val, &pokeMapAreas); err != nil {
			return PokeMapAreas{}, err
		}
		return pokeMapAreas, nil
	}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return PokeMapAreas{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return PokeMapAreas{}, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeMapAreas{}, err
	}
	pokeMapAreas := PokeMapAreas{}
	if err := json.Unmarshal(data, &pokeMapAreas); err != nil {
		return PokeMapAreas{}, err
	}
	c.cache.Add(requestURL, data)
	return pokeMapAreas, nil
}

func (c *Client) GetMapArea(mapAreaName *string) (PokeMapArea, error) {
	if mapAreaName == nil {
		return PokeMapArea{}, fmt.Errorf("empty map area name allowed")
	}
	requestURL := mapAreaEndpoint + fmt.Sprintf("/%s", *mapAreaName)

	if val, ok := c.cache.Get(requestURL); ok {
		fmt.Println("Providing cached result from past", c.cache.GetTTL(), "minutes")
		pokeMapArea := PokeMapArea{}
		if err := json.Unmarshal(val, &pokeMapArea); err != nil {
			return PokeMapArea{}, err
		}
		return pokeMapArea, nil
	}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return PokeMapArea{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return PokeMapArea{}, err
	}

	if res.StatusCode != http.StatusOK {
		return PokeMapArea{}, fmt.Errorf("received the response with %v status", res.StatusCode)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeMapArea{}, err
	}
	pokeMapArea := PokeMapArea{}
	if err := json.Unmarshal(data, &pokeMapArea); err != nil {
		return PokeMapArea{}, err
	}
	c.cache.Add(requestURL, data)
	return pokeMapArea, nil
}
