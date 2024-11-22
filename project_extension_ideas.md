#### Ideas for extending the project

*From boot.dev final pokedexcli lesson*

You don't have to extend this project, but if you're planning to make this something you add to your personal portfolio, you should consider making it your own by adding to it. Here are just a few ideas.

- Update the CLI to support the "up" arrow to cycle through previous commands
- Simulate battles between pokemon
- Add more unit tests
- Refactor your code to organize it better and make it more testable
- Keep pokemon in a "party" and allow them to level up
- Allow for pokemon that are caught to evolve after a set amount of time
- Persist a user's Pokedex to disk so they can save progress between sessions
- Use the PokeAPI to make exploration more interesting. For example, rather than typing the names of areas, maybe you are given choices of areas and just type "left" or "right"
- Random encounters with wild pokemon
- Adding support for different types of balls (Pokeballs, Great Balls, Ultra Balls, etc), which have different chances of catching pokemon


### Logging

Here's a comprehensive list of where logging would be valuable in your application, organized by component and log level:

**Cache Operations (pokecache.go)**
DEBUG: When items are added to cache
DEBUG: When items are retrieved from cache
INFO: When cache cleanup occurs
DEBUG: When items are removed during cleanup

**API Client (pokeapi.go)**
DEBUG: API requests being made
INFO: Rate limiting events
ERROR: Failed API requests
DEBUG: Cache hits vs API calls

**Pokemon Operations (command_catch.go, command_inspect.go)**
INFO: Pokemon catch attempts
INFO: Successful catches
DEBUG: Base experience and catch probability calculations
ERROR: Failed Pokemon fetches

**Location Operations (command_map.go, command_explore.go)**
DEBUG: Location navigation
INFO: Area exploration starts
ERROR: Failed location/exploration fetches

**Data Persistence (pokedata.go)**
DEBUG: File operations (already implemented)
INFO: Save operations (already implemented)
ERROR: File operation failures (already implemented)
DEBUG: Data loading attempts

**REPL (repl.go)**
DEBUG: Command parsing
ERROR: Invalid commands
INFO: Program start/stop
DEBUG: Signal handling