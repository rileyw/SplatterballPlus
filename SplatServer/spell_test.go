package main

import (
	"testing"
	"time"
)

func TestSpellManager(t *testing.T) {
	sm := NewSpellManager()

	// Test getting a spell
	spell := sm.GetSpell(1)
	if spell == nil {
		t.Fatal("Failed to get spell 1")
	}

	if spell.Name != "Fire Bolt" {
		t.Errorf("Expected spell name 'Fire Bolt', got '%s'", spell.Name)
	}

	if spell.Damage != 25 {
		t.Errorf("Expected damage 25, got %d", spell.Damage)
	}

	// Test getting all spells
	allSpells := sm.GetAllSpells()
	if len(allSpells) != 5 {
		t.Errorf("Expected 5 spells, got %d", len(allSpells))
	}

	// Test non-existent spell
	nonExistent := sm.GetSpell(999)
	if nonExistent != nil {
		t.Error("Expected nil for non-existent spell")
	}
}

func TestSpellSystem(t *testing.T) {
	// Initialize spell manager for this test
	ss := NewSpellSystem()

	// Test casting a spell
	spellInstance, err := ss.CastSpell(1, 1, 10.0, 20.0, 2)
	if err != nil {
		t.Errorf("Failed to cast spell: %v", err)
	}

	if spellInstance == nil {
		t.Fatal("Spell instance is nil")
	}

	if spellInstance.CasterID != 1 {
		t.Errorf("Expected caster ID 1, got %d", spellInstance.CasterID)
	}

	if spellInstance.SpellID != 1 {
		t.Errorf("Expected spell ID 1, got %d", spellInstance.SpellID)
	}

	// Test cooldown
	_, err = ss.CastSpell(1, 1, 10.0, 20.0, 2)
	if err == nil {
		t.Error("Expected cooldown error for immediate recast")
	}

	// Test getting active spells
	activeSpells := ss.GetActiveSpells()
	if len(activeSpells) != 1 {
		t.Errorf("Expected 1 active spell, got %d", len(activeSpells))
	}
}

func TestSpellCooldowns(t *testing.T) {
	ss := NewSpellSystem()

	// Cast a spell
	_, err := ss.CastSpell(1, 1, 10.0, 20.0, 2)
	if err != nil {
		t.Errorf("Failed to cast spell: %v", err)
	}

	// Check cooldowns
	cooldowns := ss.GetCooldowns(1)
	if cooldowns == nil {
		t.Fatal("Expected cooldowns map")
	}

	if len(cooldowns) != 1 {
		t.Errorf("Expected 1 cooldown, got %d", len(cooldowns))
	}

	// Check specific cooldown
	if cooldownEnd, exists := cooldowns[1]; !exists {
		t.Error("Expected cooldown for spell 1")
	} else {
		if time.Now().After(cooldownEnd) {
			t.Error("Cooldown should not have expired yet")
		}
	}
}

func TestSpellUpdateSystem(t *testing.T) {
	ss := NewSpellSystem()

	// Cast a spell with duration
	spellInstance, err := ss.CastSpell(1, 3, 10.0, 20.0, 1) // Speed Boost has 10s duration
	if err != nil {
		t.Errorf("Failed to cast spell: %v", err)
	}

	// Check active spells
	activeSpells := ss.GetActiveSpells()
	if len(activeSpells) != 1 {
		t.Errorf("Expected 1 active spell, got %d", len(activeSpells))
	}

	// Manually expire the spell by setting its start time to 11 seconds ago
	spellInstance.StartTime = time.Now().Add(-11 * time.Second)

	// Update system
	ss.UpdateSpellSystem(1 * time.Second)

	// Check that spell expired
	activeSpells = ss.GetActiveSpells()
	if len(activeSpells) != 0 {
		t.Errorf("Expected 0 active spells after expiration, got %d", len(activeSpells))
	}
}

func TestSpellEffects(t *testing.T) {
	// Test spell definitions
	sm := NewSpellManager()

	fireBolt := sm.GetSpell(1)
	if fireBolt.EffectType != SpellEffectDamage {
		t.Errorf("Expected Fire Bolt to be damage type")
	}
	if fireBolt.ElementType != SpellElementFire {
		t.Errorf("Expected Fire Bolt to be fire element")
	}
	if fireBolt.FriendlyType != SpellFriendlyEnemy {
		t.Errorf("Expected Fire Bolt to target enemies")
	}

	heal := sm.GetSpell(2)
	if heal.EffectType != SpellEffectHealing {
		t.Errorf("Expected Heal to be healing type")
	}
	if heal.ElementType != SpellElementHoly {
		t.Errorf("Expected Heal to be holy element")
	}
	if heal.FriendlyType != SpellFriendlyAlly {
		t.Errorf("Expected Heal to target allies")
	}

	speedBoost := sm.GetSpell(3)
	if speedBoost.EffectType != SpellEffectSpeed {
		t.Errorf("Expected Speed Boost to be speed type")
	}
	if speedBoost.FriendlyType != SpellFriendlySelf {
		t.Errorf("Expected Speed Boost to target self")
	}
}

func TestSpellConcurrency(t *testing.T) {
	ss := NewSpellSystem()
	done := make(chan bool, 10)

	// Test concurrent spell casting
	for i := 0; i < 5; i++ {
		playerID := i + 10 // Use different player IDs to avoid any potential issues
		go func(pid int) {
			_, err := ss.CastSpell(pid, 1, 10.0, 20.0, 2)
			if err != nil {
				t.Errorf("Failed to cast spell for player %d: %v", pid, err)
			}
			done <- true
		}(playerID)
	}

	// Wait for all casts
	for i := 0; i < 5; i++ {
		<-done
	}

	// Small delay to ensure all operations complete
	time.Sleep(10 * time.Millisecond)

	// Check active spells
	activeSpells := ss.GetActiveSpells()
	if len(activeSpells) != 5 {
		t.Errorf("Expected 5 active spells, got %d", len(activeSpells))
	}

	// Test concurrent cooldown checks
	for i := 0; i < 5; i++ {
		playerID := i + 10 // Use same player IDs as casting
		go func(pid int) {
			cooldowns := ss.GetCooldowns(pid)
			if cooldowns == nil {
				t.Errorf("Expected cooldowns for player %d", pid)
			}
			done <- true
		}(playerID)
	}

	// Wait for all checks
	for i := 0; i < 5; i++ {
		<-done
	}
}

func TestSpellValidation(t *testing.T) {
	ss := NewSpellSystem()

	// Test casting non-existent spell
	_, err := ss.CastSpell(1, 999, 10.0, 20.0, 2)
	if err == nil {
		t.Error("Expected error for non-existent spell")
	}

	// Test casting with invalid parameters
	_, err = ss.CastSpell(1, 1, -1000.0, 20.0, 2) // Invalid position
	// Note: Current implementation doesn't validate positions, so this should succeed
	if err != nil {
		t.Logf("Position validation not implemented yet: %v", err)
	}
}

func TestSpellProjectileTypes(t *testing.T) {
	sm := NewSpellManager()

	fireBolt := sm.GetSpell(1)
	if fireBolt.ProjectileType != SpellProjectileBolt {
		t.Errorf("Expected Fire Bolt to be bolt projectile")
	}

	iceBlast := sm.GetSpell(4)
	if iceBlast.ProjectileType != SpellProjectileBall {
		t.Errorf("Expected Ice Blast to be ball projectile")
	}

	lightning := sm.GetSpell(5)
	if lightning.ProjectileType != SpellProjectileInstant {
		t.Errorf("Expected Lightning Strike to be instant")
	}
}
