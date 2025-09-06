package main

import (
	"fmt"
	"time"
)

const (
	TICK_RATE = 60 // ticks per second
)

// GameLoop runs the main game update loop
func GameLoop(gs *GameState) {
	ticker := time.NewTicker(time.Second / TICK_RATE)
	defer ticker.Stop()

	for range ticker.C {
		UpdateGameState(gs)
	}
}

// UpdateGameState updates the game state each tick
func UpdateGameState(gs *GameState) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	now := time.Now()

	// Check for disconnected players first (before updating LastSeen)
	for id, player := range gs.Players {
		if now.Sub(player.LastSeen) > 30*time.Second {
			fmt.Printf("Player %s (ID: %d) timed out\n", player.Name, id)
			// Save before removing (ignore errors if database is closed during shutdown)
			if err := SavePlayer(player); err != nil {
				fmt.Printf("Failed to save timed out player %d: %v\n", id, err)
			}
			delete(gs.Players, id)
		}
	}

	// Update player positions, health, etc. (only for remaining players)
	for _, player := range gs.Players {
		// Example: Simple movement or health regeneration
		player.Health = min(100, player.Health+1) // Regenerate health
		player.LastSeen = time.Now()

		// Periodic save (every 10 seconds)
		if time.Since(player.LastSeen) > 10*time.Second {
			if err := SavePlayer(player); err != nil {
				fmt.Printf("Failed to save player %d: %v\n", player.ID, err)
			}
		}
	}

	// Update arenas
	UpdateArenas(gs)
}

// UpdateArenas updates all arenas in the game state
func UpdateArenas(gs *GameState) {
	gs.ArenaManager.mu.RLock()
	arenas := make([]*Arena, 0, len(gs.ArenaManager.Arenas))
	for _, arena := range gs.ArenaManager.Arenas {
		arenas = append(arenas, arena)
	}
	gs.ArenaManager.mu.RUnlock()

	for _, arena := range arenas {
		UpdateArena(arena)
	}

	// Update spell system
	gs.SpellSystem.UpdateSpellSystem(16 * time.Millisecond) // ~60 FPS
}

// UpdateArena updates a single arena
func UpdateArena(arena *Arena) {
	arena.mu.Lock()
	defer arena.mu.Unlock()

	// Check if arena should start (minimum players, etc.)
	if arena.State == ArenaStateWaiting && arena.GetPlayerCount() >= 2 {
		arena.StartArena()
		fmt.Printf("Arena %s started with %d players\n", arena.Name, arena.GetPlayerCount())
	}

	// Update arena logic based on state
	switch arena.State {
	case ArenaStateActive:
		// Handle active arena gameplay
		// TODO: Add collision detection, projectile updates, etc.
	case ArenaStateEnded:
		// Handle arena end logic
		// TODO: Calculate winners, distribute rewards, etc.
	}
}

// min is a helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
