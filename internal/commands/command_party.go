package commands

import (
	"errors"
	"fmt"

	"github.com/sakuffo/pokedexcli/internal/pokedata"
)

func CommandParty(cfg *pokedata.Config, args ...string) error {
	cfg.Logger.Debug("Executing 'party' command")
	if len(args) == 0 {
		return CommandPartyList(cfg, args...)
	}

	subcommand := args[0]

	switch subcommand {
	case "list":
		return CommandPartyList(cfg, args[1:]...)
	case "inspect":
		if len(args) != 2 {
			cfg.Logger.Error("Inspect command called without a nickname")
			return errors.New("inspect requires a nickname")
		}
		return CommandPartyInspect(cfg, args[1])
	case "remove":
		return CommandPartyRemove(cfg, args[1])
	default:
		cfg.Logger.Error("Unknown party subcommand: %s", subcommand)
		return errors.New("unknown party subcommand")
	}
}

func CommandPartyList(cfg *pokedata.Config, args ...string) error {
	members := cfg.Party.ListMembers()
	if len(members) == 0 {
		cfg.Logger.Info("No party members found")
		fmt.Println("No party members found")
		return nil
	}

	cfg.Logger.Info("Listing party members")
	fmt.Println("Party Members:")
	for _, pokemon := range members {
		fmt.Printf(" - Name: %s | Level: %d | XP: %d | Species: %s\n", pokemon.Nickname, pokemon.Level, pokemon.Experience, pokemon.BasePokemon.Species.Name)
	}
	return nil
}

func CommandPartyInspect(cfg *pokedata.Config, nickname string) error {
	cfg.Logger.Debug("Inspecting party member: %s", nickname)
	fmt.Printf("Inspecting party member: %s\n", nickname)
	fmt.Println("--------------Inspecting------------------")
	pokemon, err := cfg.Party.GetMember(nickname)
	if err != nil {
		cfg.Logger.Error("Party member not found: %s", nickname)
		fmt.Println("Party member not found: %s", nickname)
		return err
	}

	cfg.Logger.Info("Displaying details for party member: %s", nickname)
	fmt.Printf("Name: %s\n", pokemon.Nickname)
	fmt.Printf("Level: \033[32m%d\033[0m\n", pokemon.Level)
	fmt.Printf("Experience: \033[32m%d\033[0m\n", pokemon.Experience)
	fmt.Printf("Species: \033[32m%s\033[0m\n", pokemon.BasePokemon.Species.Name)
	fmt.Printf("Height: \033[32m%d\033[0m\n", pokemon.BasePokemon.Height)
	fmt.Printf("Weight: \033[32m%d\033[0m\n", pokemon.BasePokemon.Weight)

	fmt.Printf("%s's Stats:\n", pokemon.Nickname)
	fmt.Printf("  - HP: \033[32m%d\033[0m\n", pokemon.CurrentStats.HP)
	fmt.Printf("  - Attack: \033[32m%d\033[0m\n", pokemon.CurrentStats.Attack)
	fmt.Printf("  - Defense: \033[32m%d\033[0m\n", pokemon.CurrentStats.Defense)
	fmt.Printf("  - Special Attack: \033[32m%d\033[0m\n", pokemon.CurrentStats.SpecialAttack)
	fmt.Printf("  - Special Defense: \033[32m%d\033[0m\n", pokemon.CurrentStats.SpecialDefense)
	fmt.Printf("  - Speed: \033[32m%d\033[0m\n", pokemon.CurrentStats.Speed)

	fmt.Println("Types:")
	for _, t := range pokemon.BasePokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}

func CommandPartyRemove(cfg *pokedata.Config, nickname string) error {
	cfg.Logger.Debug("Removing party member: %s", nickname)
	fmt.Printf("Removing party member: %s\n", nickname)
	err := cfg.Party.RemoveMember(nickname)
	if err != nil {
		cfg.Logger.Error("Failed to remove party member: %s", nickname)
		fmt.Println("Failed to remove party member: %s", nickname)
		return err
	}
	cfg.Logger.Info("Party member removed: %s", nickname)
	fmt.Println("Party member removed: %s", nickname)
	return nil
}
