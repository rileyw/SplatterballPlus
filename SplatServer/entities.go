package main

import (
	"net"
	"sync"
	"time"
)

// Player represents a connected player
type Player struct {
	ID       int
	Name     string
	X, Y     float64
	Health   int
	Conn     net.Conn
	LastSeen time.Time
}

// GameState holds the overall game world state
type GameState struct {
	Players       map[int]*Player
	ArenaManager  *ArenaManager
	SpellSystem   *SpellSystem
	mu            sync.RWMutex
}

// NewGameState initializes a new game state
func NewGameState() *GameState {
	return &GameState{
		Players:      make(map[int]*Player),
		ArenaManager: NewArenaManager(),
		SpellSystem:  NewSpellSystem(),
	}
}

// AddPlayer adds a player to the game state
func (gs *GameState) AddPlayer(p *Player) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.Players[p.ID] = p
}

// RemovePlayer removes a player from the game state
func (gs *GameState) RemovePlayer(id int) {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	delete(gs.Players, id)
}

// GetPlayer retrieves a player by ID
func (gs *GameState) GetPlayer(id int) (*Player, bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	p, exists := gs.Players[id]
	return p, exists
}
