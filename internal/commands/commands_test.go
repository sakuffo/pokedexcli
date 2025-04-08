package commands

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/sakuffo/pokedexcli/internal/cache"
	"github.com/sakuffo/pokedexcli/internal/config"
	"github.com/sakuffo/pokedexcli/internal/discovery"
	"github.com/sakuffo/pokedexcli/internal/logger"
	"github.com/sakuffo/pokedexcli/internal/party"
	"github.com/sakuffo/pokedexcli/internal/persistence"
	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

// osExit is a variable holding os.Exit, allowing it to be mocked for tests.
var osExit = os.Exit

// mockExit is used to prevent os.Exit(0) during tests
var mockExit func(code int)

func TestMain(m *testing.M) {
	// Override os.Exit for testing
	originalExit := osExit
	mockExit = func(code int) {
		// Do nothing or maybe record the exit code
		// fmt.Printf("os.Exit(%d) called\n", code)
	}
	osExit = mockExit // Replace the global function pointer

	// Run tests
	code := m.Run()

	// Restore original os.Exit
	osExit = originalExit
	os.Exit(code)
}

// setupTestConfig creates a minimal config for testing commands.
func setupTestConfig() *config.Config {
	testLogger := logger.New(logger.NONE) // Use NONE level for tests
	testCache := cache.NewCache(5*time.Minute, testLogger)
	testClient := pokeapi.NewClient(1*time.Second, testCache, testLogger)
	testPersistence, _ := persistence.NewPersistence(".test_pokedata.json") // Use a test file
	testPersistence.SetLogger(testLogger)
	defer os.Remove(".test_pokedata.json") // Clean up test file

	cfg := &config.Config{
		Logger:        testLogger,
		PokeapiClient: testClient,
		Persistence:   testPersistence,
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
		Discoveries:   discovery.NewDiscoveryTracker(),
		Party: &party.Party{
			Members: make([]*party.PartyPokemon, 0),
		},
	}
	return cfg
}

func TestCommandHelp(t *testing.T) {
	cfg := setupTestConfig()

	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := CommandHelp(cfg)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("CommandHelp returned an unexpected error: %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Basic check: ensure some expected text is present
	if !bytes.Contains(buf.Bytes(), []byte("Welcome to the Pokedex!")) {
		t.Errorf("Help output missing expected welcome message. Got:\n%s", output)
	}
	if !bytes.Contains(buf.Bytes(), []byte("exit:")) {
		t.Errorf("Help output missing 'exit' command description. Got:\n%s", output)
	}
}

func TestCommandExit(t *testing.T) {
	cfg := setupTestConfig()

	// We use mockExit (setup in TestMain) to prevent actual exit
	exitCalled := false
	mockExit = func(code int) {
		if code != 0 {
			t.Errorf("Expected exit code 0, got %d", code)
		}
		exitCalled = true
	}
	osExit = mockExit                   // Make sure the mock is set for this test run too
	defer func() { osExit = os.Exit }() // Restore original exit func after test

	err := CommandExit(cfg)

	if err != nil {
		// CommandExit itself should return nil, even if saving fails (it logs the error).
		// The os.Exit(0) should happen regardless.
		t.Fatalf("CommandExit returned an unexpected error: %v", err)
	}

	if !exitCalled {
		t.Errorf("os.Exit(0) was not called by CommandExit")
	}

	// Optional: Check if save file was created/written
	if _, err := os.Stat(".test_pokedata.json"); os.IsNotExist(err) {
		t.Errorf("Expected save file '.test_pokedata.json' to be created, but it wasn't")
	}
}

// TODO: Add tests for command_map
// TODO: Add tests for command_explore
// TODO: Add tests for command_catch
// TODO: Add tests for command_inspect
// TODO: Add tests for command_pokedex
// TODO: Add tests for command_party

// </rewritten_file>
