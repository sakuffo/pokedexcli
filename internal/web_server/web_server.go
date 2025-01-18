package web_server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sakuffo/pokedexcli/internal/commands"
	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func StartWebServer(cfg *pokedata.Config) {
	http.HandleFunc("/help", handleHelp(cfg))
	http.HandleFunc("/catch", handleCatch(cfg))
	// http.HandleFunc("/pokedex", handlePokedex(cfg))
	// http.HandleFunc("/pokedex/inspect/", handlePokedexInspect(cfg))
	// http.HandleFunc("/party", handleParty(cfg))
	// http.HandleFunc("/map/forward", handleMapForward(cfg))
	// http.HandleFunc("/map/back", handleMapBack(cfg))
	// http.HandleFunc("/explore", handleExplore(cfg))

	// Serve static files (e.g. HTML, CSS, JS) from static directory
	fs := http.FileServer(http.Dir("/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Starting web server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		cfg.Logger.Error("Failed to start web server: %v", err)
	}
}

func handleHelp(cfg *pokedata.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := commands.CommandHelp(cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Help command executed successfully"))
	}
}

func handleCatch(cfg *pokedata.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := commands.CommandCatch(cfg, request.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handlePokedex(cfg *pokedata.Config) error {
	return nil
}
func handlePokedexInspect(cfg *pokedata.Config) error {
	return nil
}
func handleParty(cfg *pokedata.Config) error {
	return nil
}
func handleMapForward(cfg *pokedata.Config) error {
	return nil
}
func handleMapBack(cfg *pokedata.Config) error {
	return nil
}
func handleExplore(cfg *pokedata.Config) error {
	return nil
}
