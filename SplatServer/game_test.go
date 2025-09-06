package main

import (
	"fmt"
	"testing"
	"time"
)

// TestGameStateConcurrency tests concurrent access to game state
func TestGameStateConcurrency(t *testing.T) {
	gs := NewGameState()
	done := make(chan bool, 100)

	// Test concurrent player additions
	for i := 0; i < 50; i++ {
		go func(id int) {
			player := &Player{
				ID:     id,
				Name:   fmt.Sprintf("Player%d", id),
				X:      float64(id),
				Y:      float64(id * 2),
				Health: 100,
			}
			gs.AddPlayer(player)
			done <- true
		}(i)
	}

	// Wait for all additions
	for i := 0; i < 50; i++ {
		<-done
	}

	// Verify all players were added
	if len(gs.Players) != 50 {
		t.Errorf("Expected 50 players, got %d", len(gs.Players))
	}

	// Test concurrent player retrievals
	for i := 0; i < 50; i++ {
		go func(id int) {
			player, exists := gs.GetPlayer(id)
			if !exists {
				t.Errorf("Player %d not found", id)
			}
			if player.ID != id {
				t.Errorf("Player ID mismatch: expected %d, got %d", id, player.ID)
			}
			done <- true
		}(i)
	}

	// Wait for all retrievals
	for i := 0; i < 50; i++ {
		<-done
	}

	// Test concurrent removals
	for i := 0; i < 25; i++ {
		go func(id int) {
			gs.RemovePlayer(id)
			done <- true
		}(i)
	}

	// Wait for removals
	for i := 0; i < 25; i++ {
		<-done
	}

	// Verify remaining players
	if len(gs.Players) != 25 {
		t.Errorf("Expected 25 players after removal, got %d", len(gs.Players))
	}
}

// TestGameLoopUpdates tests game loop update logic
func TestGameLoopUpdates(t *testing.T) {
	gs := NewGameState()

	// Add test players
	player1 := &Player{ID: 1, Name: "Player1", Health: 50, LastSeen: time.Now()}
	player2 := &Player{ID: 2, Name: "Player2", Health: 100, LastSeen: time.Now()}
	player3 := &Player{ID: 3, Name: "Player3", Health: 10, LastSeen: time.Now().Add(-40 * time.Second)} // Should be removed

	gs.AddPlayer(player1)
	gs.AddPlayer(player2)
	gs.AddPlayer(player3)

	// Run update
	UpdateGameState(gs)

	// Check health regeneration
	if player1.Health != 51 {
		t.Errorf("Player1 health should be 51, got %d", player1.Health)
	}
	if player2.Health != 100 {
		t.Errorf("Player2 health should remain at 100, got %d", player2.Health)
	}

	// Check timeout removal
	if len(gs.Players) != 2 {
		t.Errorf("Expected 2 players after timeout removal, got %d", len(gs.Players))
	}

	_, exists := gs.GetPlayer(3)
	if exists {
		t.Error("Player3 should have been removed due to timeout")
	}
}

// TestPlayerOperations tests player-related operations
func TestPlayerOperations(t *testing.T) {
	gs := NewGameState()

	// Test player creation
	player := &Player{
		ID:     1,
		Name:   "TestPlayer",
		X:      10.0,
		Y:      20.0,
		Health: 100,
	}

	gs.AddPlayer(player)

	// Test retrieval
	retrieved, exists := gs.GetPlayer(1)
	if !exists {
		t.Fatal("Player should exist")
	}

	if retrieved.Name != "TestPlayer" {
		t.Errorf("Expected name 'TestPlayer', got '%s'", retrieved.Name)
	}

	// Test non-existent player
	_, exists = gs.GetPlayer(999)
	if exists {
		t.Error("Non-existent player should not be found")
	}

	// Test removal
	gs.RemovePlayer(1)
	_, exists = gs.GetPlayer(1)
	if exists {
		t.Error("Player should be removed")
	}
}

// TestGameStateEdgeCases tests edge cases in game state management
func TestGameStateEdgeCases(t *testing.T) {
	gs := NewGameState()

	// Test adding nil player
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when adding nil player")
		}
	}()
	gs.AddPlayer(nil)

	// Test removing non-existent player
	gs.RemovePlayer(999) // Should not panic

	// Test getting non-existent player
	_, exists := gs.GetPlayer(999)
	if exists {
		t.Error("Non-existent player should not be found")
	}
}

// TestPerformance benchmarks game state operations
func BenchmarkGameStateOperations(b *testing.B) {
	gs := NewGameState()

	// Setup benchmark data
	players := make([]*Player, 1000)
	for i := 0; i < 1000; i++ {
		players[i] = &Player{
			ID:     i,
			Name:   fmt.Sprintf("Player%d", i),
			X:      float64(i),
			Y:      float64(i * 2),
			Health: 100,
		}
	}

	b.ResetTimer()

	// Benchmark additions
	b.Run("AddPlayers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gs.AddPlayer(players[i%1000])
		}
	})

	// Benchmark retrievals
	b.Run("GetPlayers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gs.GetPlayer(i % 1000)
		}
	})

	// Benchmark updates
	b.Run("UpdateGameState", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			UpdateGameState(gs)
		}
	})
}
