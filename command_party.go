package main

import (
	"errors"
	"fmt"
)

func commandParty(cfg *config, args ...string) error {
	cfg.logger.Debug("Executing 'party' command")
	if len(args) == 0 {
		return commandPartyList(cfg, args...)
	}

	subcommand := args[0]

	switch subcommand {
	case "list":
		return commandPartyList(cfg, args[1:]...)
	case "inspect":
		if len(args) != 2 {
			cfg.logger.Error("Inspect command called without a nickname")
			return errors.New("inspect requires a nickname")
		}
		return commandPartyInspect(cfg, args[1])
	case "remove":
		return commandPartyRemove(cfg, args[1])
	default:
		cfg.logger.Error("Unknown party subcommand: %s", subcommand)
		return errors.New("Unknown party subcommand")
	}
}

func commandPartyList(cfg *config, args ...string) error {
	members := cfg.party.ListMembers()
	if len(members) == 0 {
		cfg.logger.Info("No party members found")
		fmt.Println("No party members found")
		return nil
	}

	cfg.logger.Info("Listing party members")
	fmt.Println("Party Members:")
	for _, pokemon := range members {
		fmt.Printf(" - Name: %s | Level: %d | XP: %d | Species: %s\n", pokemon.Nickname, pokemon.Level, pokemon.Experience, pokemon.BasePokemon.Species.Name)
	}
	return nil
}

func commandPartyInspect(cfg *config, nickname string) error {
	cfg.logger.Debug("Inspecting party member: %s", nickname)
	fmt.Printf("Inspecting party member: %s\n", nickname)
	fmt.Println("--------------Inspecting------------------")
	pokemon, err := cfg.party.GetMember(nickname)
	if err != nil {
		cfg.logger.Error("Party member not found: %s", nickname)
		fmt.Println("Party member not found: %s", nickname)
		return err
	}

	cfg.logger.Info("Displaying details for party member: %s", nickname)
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

func commandPartyRemove(cfg *config, nickname string) error {
	cfg.logger.Debug("Removing party member: %s", nickname)
	fmt.Printf("Removing party member: %s\n", nickname)
	err := cfg.party.RemoveMember(nickname)
	if err != nil {
		cfg.logger.Error("Failed to remove party member: %s", nickname)
		fmt.Println("Failed to remove party member: %s", nickname)
		return err
	}
	cfg.logger.Info("Party member removed: %s", nickname)
	fmt.Println("Party member removed: %s", nickname)
	return nil
}
