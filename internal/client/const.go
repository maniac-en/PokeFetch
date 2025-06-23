package client

const (
	LIMIT           string = "20"
	baseURL         string = "https://pokeapi.co/api"
	apiVersion      string = "/v2"
	mapAreaEndpoint string = baseURL + apiVersion + "/location-area"
	pokemonEndpoint string = baseURL + apiVersion + "/pokemon"
)
