# Project Progression and Code Review (As of 2024-08-01)

This document tracks the development progress and provides a code review based on Staff+ engineering principles.

## **1. Completed Components and Functionalities**

*(Existing content remains largely accurate, representing the functional state before major refactoring)*

- **CLI Commands**: `help`, `exit`, `map`, `mapb`, `explore`, `catch`, `inspect`, `pokedex`, `party`.
- **Data Persistence**: JSON storage (`.pokedexclidata/pokedata.json`), basic Load/Save with mutex protection.
- **Logging**: Configurable levels, file/stdout output.
- **API Integration**: PokeAPI client with caching.
- **Caching Mechanism**: Timed cache for API responses.
- **User Party Management**: Basic listing, inspecting, removing via commands.
- **Pokemon Inspection**: Detailed view of caught Pokemon.
- **REPL**: Interactive shell with command dispatching and graceful exit handling.
- **Configuration**: Central `Config` struct holds application state.
- **Dependency Management**: Go Modules (`go.mod`, `go.sum`).
- **Git Ignored Files**: Standard `.gitignore`.

## **2. Code Review Summary (Staff+ Perspective)**

**Strengths:**

*   Good modularity via `internal` packages (`pokeapi`, `cache`, `logger`, `repl`, etc.).
*   Clear CLI structure with distinct commands.
*   Functional persistence and caching mechanisms.
*   Concurrency safety in cache and persistence using mutexes.
*   Graceful shutdown (`Ctrl+C`) saves data.

**Areas for Improvement & Refactoring:**

*   **Overloaded Package (`internal/pokedata`):** Handles persistence, configuration, data structure definition, and discovery logic. Violates Single Responsibility Principle.
*   **Complex Initialization (`pokedata.New`):** Centralizes creation of nearly all components, leading to tight coupling and making testing difficult. Use of `log.Fatal` hinders graceful error handling in `main`.
*   **Dependency Management:** Monolithic `Config` struct obscures true dependencies of commands. Dependency Injection can be improved.
*   **Error Handling:** Inconsistent (some returns, some prints). Lack of error wrapping loses context. Commands should return errors to REPL for central handling.
*   **Domain Logic Encapsulation:** `internal/party` logic (add, remove, size checks) is scattered in commands instead of being encapsulated in `Party` methods.
*   **Discovery Serialization:** `DiscoveryTracker`'s use of a struct key (`LocationPokemon`) complicates JSON serialization, requiring `ToMap`/`FromMap`. Alternatives could be considered.
*   **Context Propagation:** Lack of `context.Context` for managing timeouts/cancellation in API calls or potentially long operations.
*   **Testing:** Low unit test coverage, especially for critical paths like commands, persistence, and domain logic. Need for more mocking via interfaces.

## **3. Proposed Refactoring Plan & Next Steps**

The following refactoring steps are proposed to improve maintainability, testability, and robustness:

1.  **Deconstruct `internal/pokedata`:**
    *   Create `internal/persistence` (for `Persistence`, `Data` struct, Load/Save).
    *   Create `internal/discovery` (for `DiscoveryTracker`, `LocationPokemon`).
    *   Create `internal/config` or `internal/app` (for `Config`/`AppState` struct - *only* runtime state).
2.  **Centralize Initialization in `main.go`:**
    *   Initialize logger, persistence, cache, API client, discovery, party manager individually in `main`.
    *   Inject *only necessary* dependencies into commands and the REPL loop.
    *   Remove the complex `pokedata.New` function.
    *   Handle initialization errors gracefully in `main`.
3.  **Standardize Error Handling:**
    *   Implement consistent error wrapping (`fmt.Errorf("...: %w", err)`).
    *   Ensure command callbacks return errors.
    *   Handle command errors centrally in the REPL.
4.  **Encapsulate Party Logic:**
    *   Add methods (`AddMember`, `RemoveMember`, `IsFull`, etc.) to the `Party` type in `internal/party`.
    *   Refactor party commands to use these methods.
5.  **Introduce `context.Context`:**
    *   Add `context.Context` to `pokeapi.Client` methods and propagate from callers.
6.  **Increase Test Coverage:**
    *   Write unit tests for persistence, discovery, party logic, and commands.
    *   Define interfaces where needed (e.g., `Persistence`, `APIClient`) to facilitate mocking.

*(Previous sections on potential areas for completion/extension can be merged or updated based on this refactoring plan)*
