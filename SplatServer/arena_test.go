package main

import (
	"testing"
)

func TestArenaManager(t *testing.T) {
	am := NewArenaManager()

	// Test creating an arena
	arena := am.CreateArena(1, "Test Arena", 4, 1)
	if arena == nil {
		t.Fatal("Failed to create arena")
	}

	if arena.ID != 1 {
		t.Errorf("Expected arena ID 1, got %d", arena.ID)
	}

	if arena.Name != "Test Arena" {
		t.Errorf("Expected arena name 'Test Arena', got '%s'", arena.Name)
	}

	if arena.MaxPlayers != 4 {
		t.Errorf("Expected max players 4, got %d", arena.MaxPlayers)
	}

	// Test getting an arena
	retrieved := am.GetArena(1)
	if retrieved == nil {
		t.Fatal("Failed to retrieve arena")
	}

	if retrieved.ID != arena.ID {
		t.Errorf("Retrieved arena ID mismatch")
	}
}

func TestArenaPlayerManagement(t *testing.T) {
	arena := &Arena{
		ID:         1,
		Name:       "Test Arena",
		MaxPlayers: 2,
		Players:    make(map[int]*ArenaPlayer),
		State:      ArenaStateWaiting,
	}

	// Test adding players
	err := arena.AddPlayer(1, TeamChaos)
	if err != nil {
		t.Errorf("Failed to add player 1: %v", err)
	}

	err = arena.AddPlayer(2, TeamBalance)
	if err != nil {
		t.Errorf("Failed to add player 2: %v", err)
	}

	if arena.GetPlayerCount() != 2 {
		t.Errorf("Expected 2 players, got %d", arena.GetPlayerCount())
	}

	// Test arena full
	err = arena.AddPlayer(3, TeamOrder)
	if err == nil {
		t.Error("Expected error when adding player to full arena")
	}

	// Test getting player
	player := arena.GetPlayer(1)
	if player == nil {
		t.Fatal("Failed to get player 1")
	}

	if player.PlayerID != 1 {
		t.Errorf("Expected player ID 1, got %d", player.PlayerID)
	}

	if player.Team != TeamChaos {
		t.Errorf("Expected team Chaos, got %d", player.Team)
	}

	// Test removing player
	arena.RemovePlayer(1)
	if arena.GetPlayerCount() != 1 {
		t.Errorf("Expected 1 player after removal, got %d", arena.GetPlayerCount())
	}

	if arena.GetPlayer(1) != nil {
		t.Error("Player 1 should be removed")
	}
}

func TestArenaStateManagement(t *testing.T) {
	arena := &Arena{
		ID:         1,
		Name:       "Test Arena",
		MaxPlayers: 4,
		Players:    make(map[int]*ArenaPlayer),
		State:      ArenaStateWaiting,
	}

	// Test initial state
	if arena.State != ArenaStateWaiting {
		t.Errorf("Expected initial state Waiting, got %d", arena.State)
	}

	// Add players and start arena
	arena.AddPlayer(1, TeamChaos)
	arena.AddPlayer(2, TeamBalance)

	// Manually start arena (normally done by game loop)
	arena.StartArena()

	if arena.State != ArenaStateActive {
		t.Errorf("Expected state Active after start, got %d", arena.State)
	}

	if arena.StartTime.IsZero() {
		t.Error("Start time should be set")
	}

	// End arena
	arena.EndArena()

	if arena.State != ArenaStateEnded {
		t.Errorf("Expected state Ended after end, got %d", arena.State)
	}

	if arena.EndTime.IsZero() {
		t.Error("End time should be set")
	}
}

func TestArenaPositionUpdates(t *testing.T) {
	arena := &Arena{
		ID:         1,
		Name:       "Test Arena",
		MaxPlayers: 4,
		Players:    make(map[int]*ArenaPlayer),
		State:      ArenaStateWaiting,
	}

	arena.AddPlayer(1, TeamChaos)

	// Test position update
	arena.UpdatePlayerPosition(1, 10.5, 20.3)

	player := arena.GetPlayer(1)
	if player == nil {
		t.Fatal("Player not found")
	}

	if player.X != 10.5 {
		t.Errorf("Expected X 10.5, got %f", player.X)
	}

	if player.Y != 20.3 {
		t.Errorf("Expected Y 20.3, got %f", player.Y)
	}
}

func TestArenaConcurrency(t *testing.T) {
	arena := &Arena{
		ID:         1,
		Name:       "Test Arena",
		MaxPlayers: 10,
		Players:    make(map[int]*ArenaPlayer),
		State:      ArenaStateWaiting,
	}

	// Test concurrent player additions
	done := make(chan bool, 10)

	for i := 0; i < 5; i++ {
		go func(id int) {
			arena.AddPlayer(id, TeamChaos)
			done <- true
		}(i)
	}

	// Wait for all additions
	for i := 0; i < 5; i++ {
		<-done
	}

	if arena.GetPlayerCount() != 5 {
		t.Errorf("Expected 5 players, got %d", arena.GetPlayerCount())
	}

	// Test concurrent position updates
	for i := 0; i < 5; i++ {
		go func(id int) {
			arena.UpdatePlayerPosition(id, float64(id*10), float64(id*20))
			done <- true
		}(i)
	}

	// Wait for all updates
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify positions
	for i := 0; i < 5; i++ {
		player := arena.GetPlayer(i)
		if player == nil {
			t.Errorf("Player %d not found", i)
			continue
		}
		expectedX := float64(i * 10)
		expectedY := float64(i * 20)
		if player.X != expectedX || player.Y != expectedY {
			t.Errorf("Player %d position mismatch: expected (%.1f, %.1f), got (%.1f, %.1f)",
				i, expectedX, expectedY, player.X, player.Y)
		}
	}
}
