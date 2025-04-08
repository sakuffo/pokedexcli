package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sakuffo/pokedexcli/internal/cache"
	"github.com/sakuffo/pokedexcli/internal/logger"
)

// Helper function to create a mock server
func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)

	// Configure logger for testing (set to NONE to discard output)
	testLogger := logger.New(logger.NONE) // Use logger.New with desired level
	// Optional: Set writer to io.Discard if you want to be explicit
	// testLogger.SetWriter(io.Discard)

	// Configure cache for testing
	testCache := cache.NewCache(5*time.Minute, testLogger) // Correct parameters for NewCache

	// Create a client pointing to the mock server
	client := NewClient(5*time.Second, testCache, testLogger)
	originalBaseURL := baseURL
	baseURL = server.URL // Override baseURL for testing
	t.Cleanup(func() {   // Use t.Cleanup to restore baseURL after test
		baseURL = originalBaseURL
	})

	return server, &client
}

func TestListLocations(t *testing.T) {
	t.Run("Successful fetch - no pagination", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/location-area" {
				t.Errorf("Expected path /location-area, got %s", r.URL.Path)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			expectedResp := Locations{
				Count:    1,
				Next:     nil,
				Previous: nil,
				Results: []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					{Name: "test-location", URL: "http://example.com/test"},
				},
			}
			json.NewEncoder(w).Encode(expectedResp)
		})

		server, client := setupMockServer(t, handler)
		defer server.Close()

		locations, err := client.ListLocations(nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if locations.Count != 1 {
			t.Errorf("Expected count 1, got %d", locations.Count)
		}
		if len(locations.Results) != 1 || locations.Results[0].Name != "test-location" {
			t.Errorf("Unexpected results: %+v", locations.Results)
		}
		// Add more assertions as needed (cache check, etc.)
	})

	t.Run("Cache hit", func(t *testing.T) {
		// Setup initial call to populate cache
		firstCallHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			expectedResp := Locations{
				Count: 1,
				Results: []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				}{
					{Name: "cached-location", URL: "http://example.com/cached"},
				},
			}
			json.NewEncoder(w).Encode(expectedResp)
		})

		server, client := setupMockServer(t, firstCallHandler)
		defer server.Close()

		_, err := client.ListLocations(nil) // First call
		if err != nil {
			t.Fatalf("First call failed: %v", err)
		}

		// Setup handler that should NOT be called if cache works
		secondCallHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("Server should not be called on cache hit")
		})
		server.Config.Handler = secondCallHandler // Update server handler

		locations, err := client.ListLocations(nil) // Second call (should hit cache)
		if err != nil {
			t.Fatalf("Expected no error on cache hit, got %v", err)
		}

		if locations.Count != 1 || len(locations.Results) != 1 || locations.Results[0].Name != "cached-location" {
			t.Errorf("Unexpected results from cache: %+v", locations.Results)
		}
	})

	// TODO: Add tests for pagination (Next/Previous URLs)
	// TODO: Add tests for API errors (e.g., 500 status)
	// TODO: Add tests for unmarshalling errors
}

// TODO: Add TestFetchAreaPokemon
// TODO: Add TestFetchPokemonSpecies
// TODO: Add TestFetchPokemon

// TODO: Add TestFetchAreaPokemon
// TODO: Add TestFetchPokemonSpecies
// TODO: Add TestFetchPokemon
