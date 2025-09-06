package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// SpellEffectType defines the type of effect a spell has
type SpellEffectType int

const (
	SpellEffectNone SpellEffectType = iota
	SpellEffectDamage
	SpellEffectHealing
	SpellEffectSpeed
	SpellEffectShield
	SpellEffectSlow
	SpellEffectStun
)

// SpellElementType defines the elemental type of a spell
type SpellElementType int

const (
	SpellElementNone SpellElementType = iota
	SpellElementFire
	SpellElementCold
	SpellElementLight
	SpellElementVoid
	SpellElementHoly
	SpellElementEarth
	SpellElementNature
	SpellElementAir
	SpellElementMana
)

// SpellFriendlyType defines whether a spell affects friends or enemies
type SpellFriendlyType int

const (
	SpellFriendlyEnemy SpellFriendlyType = iota
	SpellFriendlyAlly
	SpellFriendlySelf
)

// SpellProjectileType defines the projectile behavior
type SpellProjectileType int

const (
	SpellProjectileInstant SpellProjectileType = iota
	SpellProjectileBolt
	SpellProjectileBall
	SpellProjectileWave
)

// Spell represents a spell definition
type Spell struct {
	ID           int
	Name         string
	Description  string
	EffectType   SpellEffectType
	ElementType  SpellElementType
	FriendlyType SpellFriendlyType
	ProjectileType SpellProjectileType
	Damage       int
	Healing      int
	Duration     time.Duration
	Cooldown     time.Duration
	Range        float64
	Speed        float64
}

// SpellManager manages all spells
type SpellManager struct {
	Spells map[int]*Spell
	mu     sync.RWMutex
}

// NewSpellManager creates a new spell manager
func NewSpellManager() *SpellManager {
	sm := &SpellManager{
		Spells: make(map[int]*Spell),
	}
	sm.initializeSpells()
	return sm
}

// initializeSpells sets up the basic spell definitions
func (sm *SpellManager) initializeSpells() {
	spells := []*Spell{
		{
			ID:           1,
			Name:         "Fire Bolt",
			Description:  "Launches a bolt of fire",
			EffectType:   SpellEffectDamage,
			ElementType:  SpellElementFire,
			FriendlyType: SpellFriendlyEnemy,
			ProjectileType: SpellProjectileBolt,
			Damage:       25,
			Duration:     2 * time.Second, // Allow time for projectile to reach target
			Range:        100.0,
			Speed:        200.0,
			Cooldown:     1 * time.Second,
		},
		{
			ID:           2,
			Name:         "Heal",
			Description:  "Restores health to an ally",
			EffectType:   SpellEffectHealing,
			ElementType:  SpellElementHoly,
			FriendlyType: SpellFriendlyAlly,
			ProjectileType: SpellProjectileInstant,
			Healing:      30,
			Range:        50.0,
			Cooldown:     3 * time.Second,
		},
		{
			ID:           3,
			Name:         "Speed Boost",
			Description:  "Increases movement speed",
			EffectType:   SpellEffectSpeed,
			ElementType:  SpellElementAir,
			FriendlyType: SpellFriendlySelf,
			ProjectileType: SpellProjectileInstant,
			Duration:     10 * time.Second,
			Range:        0.0,
			Cooldown:     15 * time.Second,
		},
		{
			ID:           4,
			Name:         "Ice Blast",
			Description:  "Freezes and damages enemies",
			EffectType:   SpellEffectSlow,
			ElementType:  SpellElementCold,
			FriendlyType: SpellFriendlyEnemy,
			ProjectileType: SpellProjectileBall,
			Damage:       20,
			Duration:     3 * time.Second,
			Range:        80.0,
			Speed:        150.0,
			Cooldown:     2 * time.Second,
		},
		{
			ID:           5,
			Name:         "Lightning Strike",
			Description:  "Instant lightning damage",
			EffectType:   SpellEffectDamage,
			ElementType:  SpellElementLight,
			FriendlyType: SpellFriendlyEnemy,
			ProjectileType: SpellProjectileInstant,
			Damage:       40,
			Range:        60.0,
			Cooldown:     4 * time.Second,
		},
	}

	for _, spell := range spells {
		sm.Spells[spell.ID] = spell
	}
}

// GetSpell gets a spell by ID
func (sm *SpellManager) GetSpell(id int) *Spell {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.Spells[id]
}

// GetAllSpells returns all available spells
func (sm *SpellManager) GetAllSpells() []*Spell {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	spells := make([]*Spell, 0, len(sm.Spells))
	for _, spell := range sm.Spells {
		spells = append(spells, spell)
	}
	return spells
}

// SpellInstance represents an active spell cast
type SpellInstance struct {
	ID         int64
	SpellID    int
	CasterID   int
	TargetID   int
	X, Y       float64
	VelocityX  float64
	VelocityY  float64
	StartTime  time.Time
	Duration   time.Duration
}

// SpellSystem manages active spells and effects
type SpellSystem struct {
	SpellManager    *SpellManager
	ActiveSpells    map[int64]*SpellInstance
	CasterCooldowns map[int]map[int]time.Time // playerID -> spellID -> cooldownEnd
	mu              sync.RWMutex
}

// NewSpellSystem creates a new spell system
func NewSpellSystem() *SpellSystem {
	return &SpellSystem{
		SpellManager:    NewSpellManager(),
		ActiveSpells:    make(map[int64]*SpellInstance),
		CasterCooldowns: make(map[int]map[int]time.Time),
	}
}

// CastSpell attempts to cast a spell
func (ss *SpellSystem) CastSpell(casterID int, spellID int, targetX, targetY float64, targetID int) (*SpellInstance, error) {
	spell := ss.SpellManager.GetSpell(spellID)
	if spell == nil {
		return nil, fmt.Errorf("spell %d not found", spellID)
	}

	// Check cooldown
	if !ss.canCastSpell(casterID, spellID) {
		return nil, fmt.Errorf("spell %d is on cooldown", spellID)
	}

	// Create spell instance
	instance := &SpellInstance{
		ID:        generateSpellID(),
		SpellID:   spellID,
		CasterID:  casterID,
		TargetID:  targetID,
		X:         targetX,
		Y:         targetY,
		StartTime: time.Now(),
		Duration:  spell.Duration,
	}

	// Set cooldown
	ss.setCooldown(casterID, spellID, spell.Cooldown)

	// Add to active spells
	ss.mu.Lock()
	ss.ActiveSpells[instance.ID] = instance
	ss.mu.Unlock()

	return instance, nil
}

// canCastSpell checks if a player can cast a spell
func (ss *SpellSystem) canCastSpell(playerID, spellID int) bool {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if playerCooldowns, exists := ss.CasterCooldowns[playerID]; exists {
		if cooldownEnd, hasCooldown := playerCooldowns[spellID]; hasCooldown {
			return time.Now().After(cooldownEnd)
		}
	}
	return true
}

// setCooldown sets a cooldown for a spell
func (ss *SpellSystem) setCooldown(playerID, spellID int, duration time.Duration) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.CasterCooldowns[playerID] == nil {
		ss.CasterCooldowns[playerID] = make(map[int]time.Time)
	}
	ss.CasterCooldowns[playerID][spellID] = time.Now().Add(duration)
}

// UpdateSpellSystem updates all active spells
func (ss *SpellSystem) UpdateSpellSystem(deltaTime time.Duration) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	now := time.Now()
	for id, spell := range ss.ActiveSpells {
		// Update spell position/movement
		spellDef := ss.SpellManager.GetSpell(spell.SpellID)
		if spellDef != nil && spellDef.Speed > 0 {
			// Move projectile (simplified - would need proper direction calculation)
			spell.X += spell.VelocityX * deltaTime.Seconds()
			spell.Y += spell.VelocityY * deltaTime.Seconds()
		}

		// Check if spell has expired
		if spell.Duration > 0 && now.Sub(spell.StartTime) > spell.Duration {
			delete(ss.ActiveSpells, id)
		}
	}
}

// GetActiveSpells returns all active spells
func (ss *SpellSystem) GetActiveSpells() []*SpellInstance {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	spells := make([]*SpellInstance, 0, len(ss.ActiveSpells))
	for _, spell := range ss.ActiveSpells {
		spells = append(spells, spell)
	}
	return spells
}

// RemoveSpell removes a spell instance
func (ss *SpellSystem) RemoveSpell(spellID int64) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	delete(ss.ActiveSpells, spellID)
}

// GetCooldowns returns current cooldowns for a player
func (ss *SpellSystem) GetCooldowns(playerID int) map[int]time.Time {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	if cooldowns, exists := ss.CasterCooldowns[playerID]; exists {
		// Return a copy
		result := make(map[int]time.Time)
		for spellID, endTime := range cooldowns {
			result[spellID] = endTime
		}
		return result
	}
	return nil
}

// Global spell manager instance
var spellManager *SpellManager
var spellSystem *SpellSystem

// InitializeSpellSystem initializes the global spell system
func InitializeSpellSystem() {
	spellManager = NewSpellManager()
	spellSystem = NewSpellSystem()
	fmt.Println("SplatServer: Spell system initialized with", len(spellManager.Spells), "spells")
}

var spellIDCounter int64

// generateSpellID generates a unique spell instance ID
func generateSpellID() int64 {
	return atomic.AddInt64(&spellIDCounter, 1) + time.Now().UnixNano()
}
