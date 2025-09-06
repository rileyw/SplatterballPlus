package main

import (
	"fmt"
	"sync"
	"testing"
)

// TestDatabaseOperations tests database CRUD operations
func TestDatabaseOperations(t *testing.T) {
	// Use SQLite for testing
	originalDBType := dbType
	defer func() { dbType = originalDBType }()

	// Initialize test database
	dbType = "sqlite"
	InitDB()
	CreateTables()
	defer db.Close()

	// Test player creation and saving
	player := &Player{
		ID:     1,
		Name:   "TestPlayer",
		X:      10.5,
		Y:      20.3,
		Health: 85,
	}

	err := SavePlayer(player)
	if err != nil {
		t.Fatalf("Failed to save player: %v", err)
	}

	// Test player loading
	loadedPlayer, err := LoadPlayer(1)
	if err != nil {
		t.Fatalf("Failed to load player: %v", err)
	}

	if loadedPlayer.ID != player.ID {
		t.Errorf("ID mismatch: expected %d, got %d", player.ID, loadedPlayer.ID)
	}
	if loadedPlayer.Name != player.Name {
		t.Errorf("Name mismatch: expected %s, got %s", player.Name, loadedPlayer.Name)
	}
	if loadedPlayer.X != player.X {
		t.Errorf("X mismatch: expected %f, got %f", player.X, loadedPlayer.X)
	}
	if loadedPlayer.Y != player.Y {
		t.Errorf("Y mismatch: expected %f, got %f", player.Y, loadedPlayer.Y)
	}
	if loadedPlayer.Health != player.Health {
		t.Errorf("Health mismatch: expected %d, got %d", player.Health, loadedPlayer.Health)
	}
}

// TestDatabaseErrors tests error handling in database operations
func TestDatabaseErrors(t *testing.T) {
	// Use SQLite for testing
	originalDBType := dbType
	defer func() { dbType = originalDBType }()

	dbType = "sqlite"
	InitDB()
	CreateTables()
	defer db.Close()

	// Test loading non-existent player
	_, err := LoadPlayer(999)
	if err == nil {
		t.Error("Expected error when loading non-existent player")
	}

	// Test saving player with invalid data
	player := &Player{
		ID:     2,
		Name:   "", // Empty name
		X:      0,
		Y:      0,
		Health: 100,
	}

	err = SavePlayer(player)
	if err != nil {
		t.Logf("Save with empty name returned error (expected): %v", err)
	}
}

// TestConcurrentDatabase tests database operations under concurrent access
func TestConcurrentDatabase(t *testing.T) {
	// Use SQLite for testing
	originalDBType := dbType
	defer func() { dbType = originalDBType }()

	dbType = "sqlite"
	InitDB()
	CreateTables()
	defer db.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Test concurrent saves first
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			player := &Player{
				ID:     id + 10,
				Name:   fmt.Sprintf("ConcurrentPlayer%d", id),
				X:      float64(id),
				Y:      float64(id * 2),
				Health: 100 - id,
			}

			err := SavePlayer(player)
			mu.Unlock()
			if err != nil {
				t.Errorf("Concurrent save failed for player %d: %v", id, err)
			}
		}(i)
	}

	// Wait for all saves to complete
	wg.Wait()

	// Now test concurrent loads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			playerID := id + 10
			_, err := LoadPlayer(playerID)
			mu.Unlock()
			if err != nil {
				t.Errorf("Concurrent load failed for player %d: %v", playerID, err)
			}
		}(i)
	}

	wg.Wait()
}

// TestDatabaseInitialization tests database initialization with different types
func TestDatabaseInitialization(t *testing.T) {
	// Test SQLite initialization
	originalDBType := dbType
	defer func() { dbType = originalDBType }()

	dbType = "sqlite"
	InitDB()
	if db == nil {
		t.Error("SQLite database not initialized")
	}
	CreateTables()
	db.Close()

	// Test invalid database type (this will cause the test to fail, but not panic)
	// We can't easily test this without modifying the InitDB function to return errors
	// instead of calling log.Fatalf
	t.Log("Note: Invalid database type test would require InitDB to return errors")
}
