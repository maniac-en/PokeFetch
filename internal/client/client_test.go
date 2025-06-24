package client

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/maniac-en/pokefetch/internal/cache"
)

// Custom RoundTripper type
type mockTransport func(*http.Request) (*http.Response, error)

func (m mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req)
}

// Helper function to create HTTP response
func createResponse(statusCode int, body string, headers map[string]string) *http.Response {
	resp := &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	for key, value := range headers {
		resp.Header.Set(key, value)
	}

	return resp
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		timeout       time.Duration
		cacheInterval time.Duration
		expectError   bool
		errorMsg      string
	}{
		{
			name:          "valid parameters",
			timeout:       5 * time.Second,
			cacheInterval: 10 * time.Minute,
			expectError:   false,
		},
		{
			name:          "zero timeout",
			timeout:       0,
			cacheInterval: 10 * time.Minute,
			expectError:   true,
			errorMsg:      "timeout must be positive",
		},
		{
			name:          "negative timeout",
			timeout:       -1 * time.Second,
			cacheInterval: 10 * time.Minute,
			expectError:   true,
			errorMsg:      "timeout must be positive",
		},
		{
			name:          "zero cache interval",
			timeout:       5 * time.Second,
			cacheInterval: 0,
			expectError:   true,
			errorMsg:      "cache interval must be positive",
		},
		{
			name:          "negative cache interval",
			timeout:       5 * time.Second,
			cacheInterval: -1 * time.Minute,
			expectError:   true,
			errorMsg:      "cache interval must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.timeout, tt.cacheInterval)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
				if client != nil {
					t.Errorf("expected nil client on error, got %v", client)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if client == nil {
					t.Errorf("expected client, got nil")
				}
			}
		})
	}
}

func TestGetMapAreas_Success(t *testing.T) {
	mockResponse := `{
		"count": 781,
		"next": "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
		"previous": null,
		"results": [
			{
				"name": "canalave-city-area",
				"url": "https://pokeapi.co/api/v2/location-area/1/"
			},
			{
				"name": "eterna-city-area",
				"url": "https://pokeapi.co/api/v2/location-area/2/"
			}
		]
	}`

	client := &Client{
		cache: cache.NewCache(3 * time.Second),
		httpClient: http.Client{
			Transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", req.Method)
				}

				expectedURL := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
				if req.URL.String() != expectedURL {
					t.Errorf("expected URL %s, got %s", expectedURL, req.URL.String())
				}

				return createResponse(http.StatusOK, mockResponse, map[string]string{
					"Content-Type": "application/json",
				}), nil
			}),
		},
	}

	result, err := client.GetMapAreas(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 781 {
		t.Errorf("expected count 781, got %d", result.Count)
	}

	if len(result.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(result.Results))
	}

	if result.Results[0].Name != "canalave-city-area" {
		t.Errorf("expected first result name 'canalave-city-area', got %s", result.Results[0].Name)
	}
}

func TestGetMapAreas_WithCustomURL(t *testing.T) {
	customURL := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	mockResponse := `{
		"count": 781,
		"next": "https://pokeapi.co/api/v2/location-area?offset=60&limit=20",
		"previous": "https://pokeapi.co/api/v2/location-area?offset=20&limit=20",
		"results": []
	}`

	client := &Client{
		cache: cache.NewCache(3 * time.Second),
		httpClient: http.Client{
			Transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				if req.URL.String() != customURL {
					t.Errorf("expected URL %s, got %s", customURL, req.URL.String())
				}

				return createResponse(http.StatusOK, mockResponse, map[string]string{
					"Content-Type": "application/json",
				}), nil
			}),
		},
	}

	result, err := client.GetMapAreas(&customURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Count != 781 {
		t.Errorf("expected count 781, got %d", result.Count)
	}
}

func TestGetMapArea_Success(t *testing.T) {
	mapAreaName := "canalave-city-area"
	mockResponse := `{
		"encounter_method_rates": [],
		"game_index": 1,
		"id": 1,
		"location": {
			"name": "canalave-city",
			"url": "https://pokeapi.co/api/v2/location/1/"
		},
		"name": "canalave-city-area",
		"names": [],
		"pokemon_encounters": []
	}`

	client := &Client{
		cache: cache.NewCache(3 * time.Second),
		httpClient: http.Client{
			Transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				expectedURL := "https://pokeapi.co/api/v2/location-area/canalave-city-area"
				if req.URL.String() != expectedURL {
					t.Errorf("expected URL %s, got %s", expectedURL, req.URL.String())
				}

				return createResponse(http.StatusOK, mockResponse, map[string]string{
					"Content-Type": "application/json",
				}), nil
			}),
		},
	}

	result, err := client.GetMapArea(&mapAreaName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}

	if result.Name != "canalave-city-area" {
		t.Errorf("expected name 'canalave-city-area', got %s", result.Name)
	}
}

func TestGetPokemon_Success(t *testing.T) {
	pokemonName := "pikachu"
	mockResponse := `{
		"height": 4,
		"name": "pikachu",
		"stats": [
			{
				"base_stat": 35,
				"effort": 0,
				"stat": {
					"name": "hp",
					"url": "https://pokeapi.co/api/v2/stat/1/"
				}
			}
		],
		"types": [
			{
				"slot": 1,
				"type": {
					"name": "electric",
					"url": "https://pokeapi.co/api/v2/type/13/"
				}
			}
		],
		"weight": 60,
		"base_experience": 112
	}`

	client := &Client{
		cache: cache.NewCache(3 * time.Second),
		httpClient: http.Client{
			Transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				expectedURL := "https://pokeapi.co/api/v2/pokemon/pikachu"
				if req.URL.String() != expectedURL {
					t.Errorf("expected URL %s, got %s", expectedURL, req.URL.String())
				}

				return createResponse(http.StatusOK, mockResponse, map[string]string{
					"Content-Type": "application/json",
				}), nil
			}),
		},
	}

	result, err := client.GetPokemon(&pokemonName)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "pikachu" {
		t.Errorf("expected name 'pikachu', got %s", result.Name)
	}

	if result.BaseExperience != 112 {
		t.Errorf("expected base experience 112, got %d", result.BaseExperience)
	}
}

func TestGetResourceFromPokeAPI_ErrorCases(t *testing.T) {
	tests := []struct {
		name          string
		url           *string
		transport     mockTransport
		expectedError string
	}{
		{
			name:          "nil URL",
			url:           nil,
			transport:     nil,
			expectedError: "request URL cannot be empty",
		},
		{
			name: "HTTP client error",
			url:  stringPtr("https://pokeapi.co/api/v2/pokemon/pikachu"),
			transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			}),
			expectedError: "network error",
		},
		{
			name: "404 not found",
			url:  stringPtr("https://pokeapi.co/api/v2/pokemon/nonexistent"),
			transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				return createResponse(http.StatusNotFound, "", map[string]string{}), nil
			}),
			expectedError: "resource not found at https://pokeapi.co/api/v2/pokemon/nonexistent",
		},
		{
			name: "500 internal server error",
			url:  stringPtr("https://pokeapi.co/api/v2/pokemon/pikachu"),
			transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				return createResponse(http.StatusInternalServerError, "", map[string]string{}), nil
			}),
			expectedError: "received the response with 500 status",
		},
		{
			name: "invalid JSON response",
			url:  stringPtr("https://pokeapi.co/api/v2/pokemon/pikachu"),
			transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				return createResponse(http.StatusOK, "invalid json", map[string]string{
					"Content-Type": "application/json",
				}), nil
			}),
			expectedError: "failed to unmarshal response:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				cache: cache.NewCache(3 * time.Second),
			}

			if tt.transport != nil {
				client.httpClient = http.Client{Transport: tt.transport}
			}

			result, err := GetResourceFromPokeAPI[Pokemon](client, tt.url)

			if err == nil {
				t.Fatalf("expected error but got none")
			}

			if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("expected error containing %q, got %q", tt.expectedError, err.Error())
			}

			// Check that zero value is returned on error
			var zero Pokemon
			if !reflect.DeepEqual(result, zero) {
				t.Errorf("expected zero value on error, got %+v", result)
			}
		})
	}
}

func TestGetResourceFromPokeAPI_CacheHit(t *testing.T) {
	url := "https://pokeapi.co/api/v2/pokemon/pikachu"
	cachedData := []byte(`{"id": 25, "name": "pikachu", "base_experience": 112}`)

	cache := cache.NewCache(3 * time.Second)
	cache.Add(url, cachedData)

	client := &Client{
		cache: cache,
		httpClient: http.Client{
			Transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				t.Errorf("HTTP request should not be made when cache hits")
				return nil, nil
			}),
		},
	}

	result, err := GetResourceFromPokeAPI[Pokemon](client, &url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 25 {
		t.Errorf("expected ID 25, got %d", result.ID)
	}

	if result.Name != "pikachu" {
		t.Errorf("expected name 'pikachu', got %s", result.Name)
	}
}

func TestGetResourceFromPokeAPI_CacheCorruptedData(t *testing.T) {
	url := "https://pokeapi.co/api/v2/pokemon/pikachu"
	corruptedData := []byte(`{"invalid": json}`)

	cache := cache.NewCache(3 * time.Second)
	cache.Add(url, corruptedData)

	client := &Client{
		cache: cache,
		httpClient: http.Client{
			Transport: mockTransport(func(req *http.Request) (*http.Response, error) {
				t.Errorf("HTTP request should not be made when cache hits")
				return nil, nil
			}),
		},
	}

	result, err := GetResourceFromPokeAPI[Pokemon](client, &url)
	if err == nil {
		t.Fatalf("expected error due to corrupted cache data")
	}

	if !strings.Contains(err.Error(), "failed to unmarshal cached data") {
		t.Errorf("expected cache unmarshal error, got %v", err)
	}

	var zero Pokemon
	if !reflect.DeepEqual(result, zero) {
		t.Errorf("expected zero value on error, got %+v", result)
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
