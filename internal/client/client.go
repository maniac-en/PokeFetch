// Package client implements data fetching for the PokeAPI (pokeapi.co)
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	LIMIT           string = "20"
	baseURL         string = "https://pokeapi.co/api"
	apiVersion      string = "/v2"
	mapAreaEndpoint string = baseURL + apiVersion + "/location-area"
)

type PokeMaps struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// type PokeMap struct {
// 	Areas []struct {
// 		Name string `json:"name"`
// 		URL  string `json:"url"`
// 	} `json:"areas"`
// 	GameIndices []struct {
// 		GameIndex  int `json:"game_index"`
// 		Generation struct {
// 			Name string `json:"name"`
// 			URL  string `json:"url"`
// 		} `json:"generation"`
// 	} `json:"game_indices"`
// 	ID    int    `json:"id"`
// 	Name  string `json:"name"`
// 	Names []struct {
// 		Language struct {
// 			Name string `json:"name"`
// 			URL  string `json:"url"`
// 		} `json:"language"`
// 		Name string `json:"name"`
// 	} `json:"names"`
// 	Regions struct {
// 		Name string `json:"name"`
// 		URL  string `json:"url"`
// 	} `json:"regions"`
// }

func GetMaps(URL string, pokeMaps *PokeMaps) error {
	var requestURL string
	if URL == "" {
		requestURL = mapAreaEndpoint + fmt.Sprintf("?limit=%s", LIMIT)
	} else {
		requestURL = URL
	}

	// Make a GET request to fetch 20 location-areas
	res, err := http.Get(requestURL)
	if err != nil {
		return err
	}
	// Unmarshal the received response into our construct
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	*pokeMaps = PokeMaps{}
	if err := json.Unmarshal(data, &pokeMaps); err != nil {
		return err
	}
	return nil
}

// func GetMap(ID int) (PokeMap, error) {
// 	// validate the input ID to construct requestURL
// 	if ID < 1 {
// 		err := fmt.Errorf("invalid ID: %d, make sure it's greater than 0", ID)
// 		return PokeMap{}, err
// 	}
// 	requestURL := mapEndpoint + fmt.Sprintf("/%d", ID)
//
// 	// Make a GET request to fetch a location a.k.a, a map
// 	res, err := http.Get(requestURL)
// 	if err != nil {
// 		return PokeMap{}, err
// 	}
//
// 	// Unmarshal the received response into our construct
// 	defer res.Body.Close()
// 	data, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return PokeMap{}, err
// 	}
// 	pokeMap := PokeMap{}
// 	if err := json.Unmarshal(data, &pokeMap); err != nil {
// 		return PokeMap{}, err
// 	}
// 	return pokeMap, nil
// }

// @@@: Might use it later!
// func GetResource[T any](URL string, resource *T) error {
// 	if URL == "" {
// 		return fmt.Errorf("error fetching resource because of an empty URL")
// 	}
//
// 	res, err := http.Get(URL)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer res.Body.Close()
// 	data, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return err
// 	}
// 	// @@@: We assume that the resource pointer sent is an empty struct
// 	if err := json.Unmarshal(data, resource); err != nil {
// 		return err
// 	}
// 	return nil
// }
