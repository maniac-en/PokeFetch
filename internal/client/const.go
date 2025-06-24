package client

const (
	LIMIT                  string = "20"
	baseURL                string = "https://pokeapi.co/api"
	apiVersion             string = "/v2"
	mapAreaEndpoint        string = baseURL + apiVersion + "/location-area"
	mapAreaDefaultEndpoint string = mapAreaEndpoint + "?offset=0&limit=20"
	pokemonEndpoint        string = baseURL + apiVersion + "/pokemon"
)
