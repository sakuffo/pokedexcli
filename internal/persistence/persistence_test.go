package persistence

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sakuffo/pokedexcli/internal/discovery"
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// setupTestPersistence creates a Persistence object for testing in a temporary directory.
func setupTestPersistence(t *testing.T) (*Persistence, string) {
	testLogger := logger.New(logger.NONE)
	tempDir := t.TempDir() // Create a temporary directory for the test
	testFilename := "test_pokedata.json"
	testFilePath := filepath.Join(tempDir, testFilename)

	// We'll directly create the Persistence object with the temp path
	// Bypassing the complex logic in NewPersistence for controlled testing
	p := &Persistence{
		filePath: testFilePath,
		logger:   testLogger, // Assign the test logger
	}
	// p.SetLogger(testLogger) // Or set it after creation if preferred

	return p, testFilePath
}

func TestSaveLoad(t *testing.T) {
	p, testFilePath := setupTestPersistence(t)

	// 1. Test saving data
	t.Run("Save Data", func(t *testing.T) {
		// Create a base Pokemon for the party member
		pikaBase := pokeapi.Pokemon{Name: "pikachu", ID: 25}

		expectedData := &Data{
			CaughtPokemon: map[string]pokeapi.Pokemon{
				"pikachu": pikaBase,
			},
			PartyMembers: []*party.PartyPokemon{
				{
					Nickname:    "Sparky",
					BasePokemon: pikaBase, // Use the BasePokemon field
					// Other fields like Level, Experience, etc., will have zero values
				},
			},
			Discoveries: discovery.NewDiscoveryTracker(), // Use an initialized tracker
		}
		// Add some discovery data for testing using the correct method
		expectedData.Discoveries.MarkDiscovered("area1", "bulbasaur")

		err := p.Save(expectedData)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Check if file exists
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Fatalf("Save did not create the file: %s", testFilePath)
		}
	})

	// 2. Test loading the saved data
	t.Run("Load Saved Data", func(t *testing.T) {
		loadedData, err := p.Load()
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		// Prepare expected data again for comparison (ensure Discoveries are comparable)
		pikaBase := pokeapi.Pokemon{Name: "pikachu", ID: 25}
		expectedData := &Data{
			CaughtPokemon: map[string]pokeapi.Pokemon{
				"pikachu": pikaBase,
			},
			PartyMembers: []*party.PartyPokemon{
				{
					Nickname:    "Sparky",
					BasePokemon: pikaBase,
				},
			},
			Discoveries: discovery.NewDiscoveryTracker(),
		}
		expectedData.Discoveries.MarkDiscovered("area1", "bulbasaur")

		if !reflect.DeepEqual(loadedData, expectedData) {
			t.Errorf("Loaded data does not match saved data.\nExpected: %+v\nGot:      %+v", expectedData, loadedData)
			// Log details for easier debugging
			t.Logf("Expected Caught: %+v", expectedData.CaughtPokemon)
			t.Logf("Got Caught: %+v", loadedData.CaughtPokemon)
			t.Logf("Expected Party: %+v", expectedData.PartyMembers)
			t.Logf("Got Party: %+v", loadedData.PartyMembers)
			t.Logf("Expected Discoveries: %+v", expectedData.Discoveries)
			t.Logf("Got Discoveries: %+v", loadedData.Discoveries)
		}
	})
}

func TestLoadNonExistent(t *testing.T) {
	p, _ := setupTestPersistence(t)

	// Ensure the file does not exist before loading
	_ = os.Remove(p.filePath) // Ignore error if it doesn't exist

	loadedData, err := p.Load()
	if err != nil {
		t.Fatalf("Load failed when file does not exist: %v", err)
	}

	// Expect empty/initialized data
	expectedData := &Data{
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
		PartyMembers:  make([]*party.PartyPokemon, 0),
		Discoveries:   discovery.NewDiscoveryTracker(),
	}

	if !reflect.DeepEqual(loadedData, expectedData) {
		t.Errorf("Loaded data from non-existent file is not empty/initialized.\nExpected: %+v\nGot:      %+v", expectedData, loadedData)
	}
}

// TODO: Test NewPersistence behavior (directory creation, permissions fallback)
// TODO: Test error handling (e.g., corrupted JSON, I/O errors during save/load)
