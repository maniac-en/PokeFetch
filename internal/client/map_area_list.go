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
		fmt.Println("Providing cached result from past 5 minutes")
		pokeMaps := PokeMapAreas{}
		if err := json.Unmarshal(val, &pokeMaps); err != nil {
			return PokeMapAreas{}, err
		}
		return pokeMaps, nil
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
	pokeMaps := PokeMapAreas{}
	if err := json.Unmarshal(data, &pokeMaps); err != nil {
		return PokeMapAreas{}, err
	}
	c.cache.Add(requestURL, data)
	return pokeMaps, nil
}
