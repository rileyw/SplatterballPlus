# Arena System Documentation

## Overview

The Arena system in SplatServer provides the core gameplay mechanics for competitive multiplayer battles. Players can join different arenas, form teams, and engage in real-time combat.

## Features

### Arena Management
- **Multiple Arenas**: Support for multiple concurrent arenas
- **Player Capacity**: Configurable maximum players per arena
- **State Management**: Waiting → Active → Ended state transitions
- **Grid Integration**: Each arena is associated with a game grid

### Team System
- **Three Teams**: Chaos, Balance, Order
- **Team Assignment**: Players choose their team when joining
- **Team Balancing**: Future implementation for automatic balancing

### Player Management
- **Arena Players**: Separate player state within arenas
- **Position Tracking**: Real-time position updates
- **Health/Score**: Player statistics within arena context

## Network Protocol

### Message Types

#### Join Arena (MsgJoinArena = 6)
```
Data: [arena_id: uint32][team: uint32]
- arena_id: ID of the arena to join
- team: Team to join (1=Chaos, 2=Balance, 3=Order)
```

#### Leave Arena (MsgLeaveArena = 7)
```
Data: [arena_id: uint32]
- arena_id: ID of the arena to leave
```

#### Arena List (MsgArenaList = 8)
```
Data: [] (empty)
Response: Text list of available arenas
```

#### Arena Update (MsgArenaUpdate = 9)
```
Data: [arena_id: uint32][x: float64][y: float64]
- arena_id: Arena ID
- x, y: Player position coordinates
```

## Usage Example

```go
// Create arena manager
arenaManager := NewArenaManager()

// Create an arena
arena := arenaManager.CreateArena(1, "Chaos Arena", 8, 1)

// Add player to arena
err := arena.AddPlayer(playerID, TeamChaos)

// Update player position
arena.UpdatePlayerPosition(playerID, 10.5, 20.3)

// Start arena when ready
if arena.GetPlayerCount() >= 2 {
    arena.StartArena()
}
```

## Game Loop Integration

The arena system integrates with the main game loop:

1. **Arena Updates**: Called every tick to update arena state
2. **Player Management**: Handles player joins/leaves
3. **State Transitions**: Manages waiting → active → ended flow
4. **Position Updates**: Processes real-time position changes

## Future Enhancements

### Planned Features
- **Projectile System**: Bolts, spells, and projectiles
- **Collision Detection**: Player and projectile collisions
- **Scoring System**: Points, kills, objectives
- **Power-ups**: Temporary abilities and bonuses
- **Arena Rulesets**: Different game modes and objectives

### Advanced Features
- **Spectator Mode**: Watch ongoing matches
- **Replay System**: Record and replay matches
- **Tournament Mode**: Bracket-based competitions
- **Custom Arenas**: Player-created arena layouts

## Testing

Run arena tests:
```bash
go test -v -run="TestArena"
```

Test with client:
```bash
cd client
go run test_client.go
```

## Configuration

Arena settings can be configured in the server initialization:

```go
// Create arenas with different settings
gs.ArenaManager.CreateArena(1, "Chaos Arena", 8, 1)    // 8 players, grid 1
gs.ArenaManager.CreateArena(2, "Balance Arena", 6, 2)  // 6 players, grid 2
gs.ArenaManager.CreateArena(3, "Order Arena", 4, 3)    // 4 players, grid 3
```

## Database Integration

Arena data is stored in the database:
- Player positions and stats
- Arena state and configuration
- Match results and history

## Performance Considerations

- **Concurrent Access**: Thread-safe arena operations
- **Memory Management**: Efficient player state storage
- **Network Optimization**: Minimal position update packets
- **Scalability**: Support for multiple concurrent arenas
