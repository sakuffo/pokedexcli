package pokeapi

// PokemonClient defines the interface for Pokemon data retrieval
type PokemonClient interface {

	// ListLocations fetches a paginated list of locations
	ListLocations(pageURL *string) (Locations, error)

	// FetchAreaPokemon fetches the Pokemon in a specific area
	FetchAreaPokemon(area *string) (Area, error)

	// FetchPokemon fetches details about a specific Pokemon
	FetchPokemon(pokemonName string) (Pokemon, error)

	// FetchPokemonSpecies species data for a specific Pokemon
	FetchPokemonSpecies(pokemonSpeciesName string) (PokemonSpecies, error)
}
