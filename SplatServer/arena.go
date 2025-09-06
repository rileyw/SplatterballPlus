package main

import (
	"fmt"
	"sync"
	"time"
)

// Arena represents a game arena
type Arena struct {
	ID          int
	Name        string
	MaxPlayers  int
	Players     map[int]*ArenaPlayer
	State       ArenaState
	StartTime   time.Time
	EndTime     time.Time
	GridID      int
	mu          sync.RWMutex
}

// ArenaState represents the current state of an arena
type ArenaState int

const (
	ArenaStateWaiting ArenaState = iota
	ArenaStateActive
	ArenaStateEnded
)

// ArenaPlayer represents a player in an arena
type ArenaPlayer struct {
	PlayerID int
	Team     Team
	X, Y     float64
	Health   int
	Score    int
}

// Team represents a team in the arena
type Team int

const (
	TeamNone Team = iota
	TeamChaos
	TeamBalance
	TeamOrder
)

// ArenaManager manages all arenas
type ArenaManager struct {
	Arenas map[int]*Arena
	mu     sync.RWMutex
}

// NewArenaManager creates a new arena manager
func NewArenaManager() *ArenaManager {
	return &ArenaManager{
		Arenas: make(map[int]*Arena),
	}
}

// CreateArena creates a new arena
func (am *ArenaManager) CreateArena(id int, name string, maxPlayers int, gridID int) *Arena {
	arena := &Arena{
		ID:         id,
		Name:       name,
		MaxPlayers: maxPlayers,
		Players:    make(map[int]*ArenaPlayer),
		State:      ArenaStateWaiting,
		GridID:     gridID,
	}

	am.mu.Lock()
	am.Arenas[id] = arena
	am.mu.Unlock()

	return arena
}

// GetArena gets an arena by ID
func (am *ArenaManager) GetArena(id int) *Arena {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.Arenas[id]
}

// AddPlayer adds a player to an arena
func (a *Arena) AddPlayer(playerID int, team Team) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if len(a.Players) >= a.MaxPlayers {
		return fmt.Errorf("arena is full")
	}

	if _, exists := a.Players[playerID]; exists {
		return fmt.Errorf("player already in arena")
	}

	arenaPlayer := &ArenaPlayer{
		PlayerID: playerID,
		Team:     team,
		X:        0,
		Y:        0,
		Health:   100,
		Score:    0,
	}

	a.Players[playerID] = arenaPlayer
	return nil
}

// RemovePlayer removes a player from an arena
func (a *Arena) RemovePlayer(playerID int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.Players, playerID)
}

// GetPlayer gets a player from the arena
func (a *Arena) GetPlayer(playerID int) *ArenaPlayer {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Players[playerID]
}

// UpdatePlayerPosition updates a player's position in the arena
func (a *Arena) UpdatePlayerPosition(playerID int, x, y float64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if player, exists := a.Players[playerID]; exists {
		player.X = x
		player.Y = y
	}
}

// StartArena starts the arena game
func (a *Arena) StartArena() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.State != ArenaStateWaiting {
		return
	}

	a.State = ArenaStateActive
	a.StartTime = time.Now()
}

// EndArena ends the arena game
func (a *Arena) EndArena() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.State != ArenaStateActive {
		return
	}

	a.State = ArenaStateEnded
	a.EndTime = time.Now()
}

// GetPlayerCount returns the number of players in the arena
func (a *Arena) GetPlayerCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.Players)
}

// IsFull checks if the arena is full
func (a *Arena) IsFull() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.Players) >= a.MaxPlayers
}
