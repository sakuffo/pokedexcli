package web_server

import (
	"fmt"
	"net/http"

	"github.com/sakuffo/pokedexcli/internal/pokeconfig"
)

func StartWebServer(cfg *pokeconfig.Config) {
	// http.HandleFunc("/help", handleHelp(cfg))
	// http.HandleFunc("/catch", handleCatch(cfg))
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

func handleHelp(cfg *pokeconfig.Config) error {
	return nil
}
func handleCatch(cfg *pokeconfig.Config) error {
	return nil
}
func handlePokedex(cfg *pokeconfig.Config) error {
	return nil
}
func handlePokedexInspect(cfg *pokeconfig.Config) error {
	return nil
}
func handleParty(cfg *pokeconfig.Config) error {
	return nil
}
func handleMapForward(cfg *pokeconfig.Config) error {
	return nil
}
func handleMapBack(cfg *pokeconfig.Config) error {
	return nil
}
func handleExplore(cfg *pokeconfig.Config) error {
	return nil
}
