package web_server

import (
	"net/http"

	"github.com/sakuffo/pokedexcli/internal/config"
)

func StartWebServer(cfg *config.Config) {
	http.HandleFunc("/help", handleHelp(cfg))
	http.HandleFunc("/catch", handleCatch(cfg))
	http.HandleFunc("/pokedex", handlePokedex(cfg))
	http.HandleFunc("/pokedex/inspect/", handlePokedexInspect(cfg))
	http.HandleFunc("/party", handleParty(cfg))
	http.HandleFunc("/map/forward", handleMapForward(cfg))
	http.HandleFunc("/map/back", handleMapBack(cfg))
	http.HandleFunc("/explore", handleExplore(cfg))

}
