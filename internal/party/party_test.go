package party

import (
	"fmt"
	"testing"

	"github.com/sakuffo/pokedexcli/internal/pokeapi"
)

func TestAddMember(t *testing.T) {
	party := &Party{Members: make([]*PartyPokemon, 0)}
	pikaBase := pokeapi.Pokemon{Name: "pikachu", ID: 25}
	pikaParty := NewPartyPokemon(pikaBase)

	// Add first member
	err := party.AddMember(pikaParty)
	if err != nil {
		t.Fatalf("Failed to add first member: %v", err)
	}
	if len(party.Members) != 1 || party.Members[0].BasePokemon.Name != "pikachu" {
		t.Errorf("Party state incorrect after adding first member: %+v", party.Members)
	}

	// Try adding the same Pokemon again
	err = party.AddMember(pikaParty) // Should fail
	if err == nil {
		t.Errorf("Expected error when adding duplicate Pokemon, got nil")
	}
	if len(party.Members) != 1 { // Length should not change
		t.Errorf("Party size changed after attempting to add duplicate")
	}

	// Fill the party
	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("pokemon-%d", i)
		tempBase := pokeapi.Pokemon{Name: name, ID: i + 100}
		tempParty := NewPartyPokemon(tempBase)
		err = party.AddMember(tempParty)
		if err != nil {
			t.Fatalf("Failed to add member %d: %v", i+2, err)
		}
	}

	if len(party.Members) != 6 {
		t.Errorf("Expected party size 6 after filling, got %d", len(party.Members))
	}

	// Try adding to a full party
	fullBase := pokeapi.Pokemon{Name: "extra", ID: 999}
	fullParty := NewPartyPokemon(fullBase)
	err = party.AddMember(fullParty) // Should fail
	if err == nil {
		t.Errorf("Expected error when adding to full party, got nil")
	}
	if len(party.Members) != 6 { // Length should not change
		t.Errorf("Party size changed after attempting to add to full party")
	}
}

func TestRemoveMember(t *testing.T) {
	party := &Party{Members: make([]*PartyPokemon, 0)}
	pikaBase := pokeapi.Pokemon{Name: "pikachu", ID: 25}
	bulbaBase := pokeapi.Pokemon{Name: "bulbasaur", ID: 1}
	pikaParty := NewPartyPokemon(pikaBase)
	bulbaParty := NewPartyPokemon(bulbaBase)

	_ = party.AddMember(pikaParty)
	_ = party.AddMember(bulbaParty)

	// Remove existing member
	err := party.RemoveMember("pikachu")
	if err != nil {
		t.Fatalf("Failed to remove 'pikachu': %v", err)
	}
	if len(party.Members) != 1 || party.Members[0].BasePokemon.Name != "bulbasaur" {
		t.Errorf("Party state incorrect after removing 'pikachu': %+v", party.Members)
	}

	// Try removing non-existent member
	err = party.RemoveMember("charmander")
	if err == nil {
		t.Errorf("Expected error when removing non-existent Pokemon, got nil")
	}
	if len(party.Members) != 1 { // Length should not change
		t.Errorf("Party size changed after attempting to remove non-existent member")
	}

	// Remove last member
	err = party.RemoveMember("bulbasaur")
	if err != nil {
		t.Fatalf("Failed to remove 'bulbasaur': %v", err)
	}
	if len(party.Members) != 0 {
		t.Errorf("Party should be empty after removing last member, size is %d", len(party.Members))
	}

	// Try removing from empty party
	err = party.RemoveMember("pikachu")
	if err == nil {
		t.Errorf("Expected error when removing from empty party, got nil")
	}
}

func TestListMembers(t *testing.T) {
	party := &Party{Members: make([]*PartyPokemon, 0)}
	pikaBase := pokeapi.Pokemon{Name: "pikachu", ID: 25}
	bulbaBase := pokeapi.Pokemon{Name: "bulbasaur", ID: 1}
	pikaParty := NewPartyPokemon(pikaBase)
	bulbaParty := NewPartyPokemon(bulbaBase)

	_ = party.AddMember(pikaParty)
	_ = party.AddMember(bulbaParty)

	members := party.ListMembers()
	if len(members) != 2 {
		t.Fatalf("Expected 2 members, got %d", len(members))
	}

	// Check if the correct members are present (order might vary)
	foundPika := false
	foundBulba := false
	for _, m := range members {
		if m.BasePokemon.Name == "pikachu" {
			foundPika = true
		} else if m.BasePokemon.Name == "bulbasaur" {
			foundBulba = true
		}
	}
	if !foundPika || !foundBulba {
		t.Errorf("Listed members did not contain expected Pokemon: %+v", members)
	}

	// Ensure the returned list is a copy (modifying it shouldn't affect original party)
	if len(members) > 0 {
		members[0].Nickname = "MODIFIED"
		if party.Members[0].Nickname == "MODIFIED" {
			t.Errorf("ListMembers returned a slice that aliases the internal party members")
		}
	}
}

func TestGetMember(t *testing.T) {
	party := &Party{Members: make([]*PartyPokemon, 0)}
	pikaBase := pokeapi.Pokemon{Name: "pikachu", ID: 25}
	pikaParty := NewPartyPokemon(pikaBase)
	_ = party.AddMember(pikaParty)

	// Get existing member
	member, found := party.GetMember("pikachu")
	if !found {
		t.Errorf("Expected to find 'pikachu', but not found")
	}
	if member == nil || member.BasePokemon.Name != "pikachu" {
		t.Errorf("GetMember returned incorrect member for 'pikachu': %+v", member)
	}

	// Get non-existent member
	_, found = party.GetMember("charmander")
	if found {
		t.Errorf("Expected not to find 'charmander', but was found")
	}
}

func TestIsFull(t *testing.T) {
	party := &Party{Members: make([]*PartyPokemon, 0)}
	if party.IsFull() {
		t.Errorf("Empty party reported as full")
	}

	// Add 5 members
	for i := 0; i < 5; i++ {
		tempBase := pokeapi.Pokemon{Name: fmt.Sprintf("poke%d", i), ID: i}
		_ = party.AddMember(NewPartyPokemon(tempBase))
	}
	if party.IsFull() {
		t.Errorf("Party with 5 members reported as full")
	}

	// Add 6th member
	sixthBase := pokeapi.Pokemon{Name: "poke6", ID: 6}
	_ = party.AddMember(NewPartyPokemon(sixthBase))
	if !party.IsFull() {
		t.Errorf("Party with 6 members not reported as full")
	}
}
