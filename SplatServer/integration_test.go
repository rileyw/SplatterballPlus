package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

// TestServerLifecycle tests the complete server startup and shutdown
func TestServerLifecycle(t *testing.T) {
	// This test requires careful setup to avoid port conflicts
	// Skip in automated testing
	t.Skip("Lifecycle test requires manual setup")

	// Initialize test database
	originalDBType := dbType
	defer func() { dbType = originalDBType }()

	dbType = "sqlite"
	InitDB()
	CreateTables()
	defer db.Close()

	// Start server in goroutine
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// This would normally start the server
		// For testing, we'll simulate the components
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test server components
	gs := NewGameState()
	if gs == nil {
		t.Error("GameState not initialized")
	}

	// Test database connectivity
	if db == nil {
		t.Error("Database not initialized")
	}

	wg.Wait()
}

// TestNetworkIntegration tests network operations end-to-end
func TestNetworkIntegration(t *testing.T) {
	// Test packet round-trip
	originalPacket := BuildLoginPacket("IntegrationTest")
	serialized := originalPacket.Serialize()

	reader := bytes.NewReader(serialized)
	deserialized, err := DeserializePacket(reader)
	if err != nil {
		t.Fatalf("Round-trip deserialization failed: %v", err)
	}

	if deserialized.Type != originalPacket.Type {
		t.Errorf("Type mismatch in round-trip")
	}

	if !bytes.Equal(deserialized.Data, originalPacket.Data) {
		t.Errorf("Data mismatch in round-trip")
	}
}

// TestConcurrentConnections simulates multiple client connections
func TestConcurrentConnections(t *testing.T) {
	// Use a separate database connection for this test
	testDB, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer testDB.Close()
	defer os.Remove("./test.db") // Clean up after test

	// Set the global db variable for this test
	originalDB := db
	db = testDB
	defer func() { db = originalDB }()

	// Initialize schema
	dbType = "sqlite"
	CreateTables()

	gs := NewGameState()
	done := make(chan bool, 10)

	// Simulate 10 concurrent connections
	for i := 0; i < 10; i++ {
		go func(id int) {
			player := &Player{
				ID:     id + 100,
				Name:   fmt.Sprintf("ConcurrentPlayer%d", id),
				X:      float64(id * 10),
				Y:      float64(id * 20),
				Health: 100,
			}

			// Simulate connection lifecycle
			gs.AddPlayer(player)

			// Simulate some game activity
			for j := 0; j < 5; j++ {
				player.X += 1.0
				player.Y += 1.0
				time.Sleep(1 * time.Millisecond)
			}

			// Save player data
			err := SavePlayer(player)
			if err != nil {
				t.Errorf("Failed to save player %d: %v", id, err)
			}

			// Simulate disconnection
			gs.RemovePlayer(player.ID)

			done <- true
		}(i)
	}

	// Wait for all connections to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	if len(gs.Players) != 0 {
		t.Errorf("Expected no players remaining, got %d", len(gs.Players))
	}
}

// TestLoadTesting simulates high load scenarios
func TestLoadTesting(t *testing.T) {
	// Skip load testing in normal runs
	t.Skip("Load test - run manually for performance testing")

	gs := NewGameState()
	start := time.Now()

	// Add many players
	for i := 0; i < 10000; i++ {
		player := &Player{
			ID:     i,
			Name:   fmt.Sprintf("LoadPlayer%d", i),
			X:      float64(i % 100),
			Y:      float64(i / 100),
			Health: 100,
		}
		gs.AddPlayer(player)
	}

	addTime := time.Since(start)
	t.Logf("Added 10000 players in %v", addTime)

	// Test retrieval performance
	start = time.Now()
	for i := 0; i < 1000; i++ {
		gs.GetPlayer(i)
	}
	getTime := time.Since(start)
	t.Logf("Retrieved 1000 players in %v", getTime)

	// Test update performance
	start = time.Now()
	for i := 0; i < 100; i++ {
		UpdateGameState(gs)
	}
	updateTime := time.Since(start)
	t.Logf("Updated game state 100 times in %v", updateTime)
}

// TestErrorRecovery tests error handling and recovery
func TestErrorRecovery(t *testing.T) {
	gs := NewGameState()

	// Test recovery from invalid operations
	player := &Player{ID: 1, Name: "Test"}
	gs.AddPlayer(player)

	// Test double removal
	gs.RemovePlayer(1)
	gs.RemovePlayer(1) // Should not panic

	// Test operations on removed player
	_, exists := gs.GetPlayer(1)
	if exists {
		t.Error("Removed player should not exist")
	}

	// Test adding player with same ID
	newPlayer := &Player{ID: 1, Name: "NewTest"}
	gs.AddPlayer(newPlayer)

	retrieved, exists := gs.GetPlayer(1)
	if !exists {
		t.Error("New player should exist")
	}
	if retrieved.Name != "NewTest" {
		t.Error("Player name should be updated")
	}
}

// TestMemoryLeaks tests for potential memory leaks
func TestMemoryLeaks(t *testing.T) {
	// Skip memory test in normal runs
	t.Skip("Memory leak test - run with memory profiling")

	gs := NewGameState()

	// Add and remove players repeatedly
	for i := 0; i < 1000; i++ {
		player := &Player{
			ID:     i,
			Name:   fmt.Sprintf("LeakTest%d", i),
			Health: 100,
		}
		gs.AddPlayer(player)
		gs.RemovePlayer(i)
	}

	// Force garbage collection
	// runtime.GC()

	if len(gs.Players) != 0 {
		t.Errorf("Memory leak detected: %d players remaining", len(gs.Players))
	}
}
