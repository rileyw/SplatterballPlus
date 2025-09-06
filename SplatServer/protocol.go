package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// MessageType defines the type of message
type MessageType byte

const (
	MsgLogin       MessageType = 1
	MsgMove        MessageType = 2
	MsgChat        MessageType = 3
	MsgLogout      MessageType = 4
	MsgPing        MessageType = 5
	MsgJoinArena   MessageType = 6
	MsgLeaveArena  MessageType = 7
	MsgArenaList   MessageType = 8
	MsgArenaUpdate MessageType = 9
	MsgCastSpell   MessageType = 10
	MsgSpellList   MessageType = 11
)

// DebugPacketCapture represents a captured unhandled packet
type DebugPacketCapture struct {
	MessageType MessageType
	Data        []byte
	PlayerID    int
	PlayerName  string
	Timestamp   string
}

// DebugPacketLogger manages captured packets
type DebugPacketLogger struct {
	CapturedPackets []DebugPacketCapture
	MaxPackets      int
	Enabled         bool
}

var debugLogger = initializeDebugLogger()

// initializeDebugLogger creates and configures the debug logger
func initializeDebugLogger() *DebugPacketLogger {
	logger := &DebugPacketLogger{
		CapturedPackets: make([]DebugPacketCapture, 0),
		MaxPackets:      100, // Default value
		Enabled:         false, // Default disabled
	}

	// Check environment variables
	if enabled := os.Getenv("DEBUG_PACKETS_ENABLED"); enabled == "true" || enabled == "1" {
		logger.Enabled = true
	}

	if maxPackets := os.Getenv("DEBUG_PACKETS_MAX"); maxPackets != "" {
		if val, err := strconv.Atoi(maxPackets); err == nil && val > 0 {
			logger.MaxPackets = val
		}
	}

	return logger
}

// Message represents a parsed message
type Message struct {
	Type MessageType
	Data []byte
}

// ParseMessage reads and parses a message from the connection
func ParseMessage(conn net.Conn) (*Message, error) {
	packet, err := DeserializePacket(conn)
	if err != nil {
		return nil, err
	}

	// Convert Packet to Message for backward compatibility
	msg := &Message{
		Type: MessageType(packet.Type),
		Data: packet.Data,
	}
	return msg, nil
}

// HandleMessage dispatches the message to the appropriate handler
func HandleMessage(msg *Message, player *Player, gs *GameState) {
	switch msg.Type {
	case MsgLogin:
		handleLogin(msg, player, gs)
	case MsgMove:
		handleMove(msg, player, gs)
	case MsgChat:
		handleChat(msg, player, gs)
	case MsgLogout:
		handleLogout(msg, player, gs)
	case MsgPing:
		handlePing(msg, player, gs)
	case MsgJoinArena:
		handleJoinArena(msg, player, gs)
	case MsgLeaveArena:
		handleLeaveArena(msg, player, gs)
	case MsgArenaList:
		handleArenaList(msg, player, gs)
	case MsgArenaUpdate:
		handleArenaUpdate(msg, player, gs)
	case MsgCastSpell:
		handleCastSpell(msg, player, gs)
	case MsgSpellList:
		handleSpellList(msg, player, gs)
	default:
		handleUnknownMessage(msg, player)
	}
}

// handleLogin processes a login message
func handleLogin(msg *Message, player *Player, gs *GameState) {
	name, err := ParseLoginPacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse login: %v\n", err)
		return
	}
	player.Name = name

	// Try to load existing player data
	if existingPlayer, err := LoadPlayer(player.ID); err == nil {
		player.Name = existingPlayer.Name
		player.X = existingPlayer.X
		player.Y = existingPlayer.Y
		player.Health = existingPlayer.Health
		fmt.Printf("Loaded existing player %s\n", player.Name)
	} else {
		// New player, save initial data
		if err := SavePlayer(player); err != nil {
			fmt.Printf("Failed to save new player: %v\n", err)
		}
	}

	fmt.Printf("Player %d logged in as %s\n", player.ID, player.Name)
	// TODO: Authenticate, load player data, etc.
}

// handleMove processes a move message
func handleMove(msg *Message, player *Player, gs *GameState) {
	x, y, err := ParseMovePacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse move: %v\n", err)
		return
	}
	player.X = x
	player.Y = y
	fmt.Printf("Player %d moved to (%.2f, %.2f)\n", player.ID, player.X, player.Y)
	// TODO: Validate move, update game state, broadcast to others
}

// handleChat processes a chat message
func handleChat(msg *Message, player *Player, gs *GameState) {
	chatMsg, err := ParseChatPacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse chat: %v\n", err)
		return
	}

	// Check for debug commands
	if handleDebugCommand(chatMsg, player) {
		return // Command was handled, don't broadcast as regular chat
	}

	fmt.Printf("Player %d said: %s\n", player.ID, chatMsg)
	// TODO: Broadcast chat to other players
}

// handleDebugCommand processes debug commands
func handleDebugCommand(message string, player *Player) bool {
	if !strings.HasPrefix(message, "/debug") {
		return false
	}

	// Check if debug logging is enabled
	if !debugLogger.Enabled {
		player.Conn.Write([]byte("Debug packet capture is disabled. Set DEBUG_PACKETS_ENABLED=true to enable.\n"))
		return true
	}

	parts := strings.Fields(message)
	if len(parts) < 2 {
		sendDebugHelp(player)
		return true
	}

	command := parts[1]
	switch command {
	case "packets":
		sendCapturedPackets(player)
	case "stats":
		sendDebugStats(player)
	case "clear":
		debugLogger.ClearCapturedPackets()
		player.Conn.Write([]byte("Debug packets cleared\n"))
	case "help":
		sendDebugHelp(player)
	default:
		player.Conn.Write([]byte("Unknown debug command. Use /debug help for available commands\n"))
	}

	return true
}

// handleLogout processes a logout message
func handleLogout(msg *Message, player *Player, gs *GameState) {
	// Save player data before logout
	if err := SavePlayer(player); err != nil {
		fmt.Printf("Failed to save player %d on logout: %v\n", player.ID, err)
	} else {
		fmt.Printf("Saved player %d data on logout\n", player.ID)
	}

	fmt.Printf("Player %d logged out\n", player.ID)
	gs.RemovePlayer(player.ID)
	player.Conn.Close()
}

// handlePing processes a ping message
func handlePing(msg *Message, player *Player, gs *GameState) {
	// Respond with pong
	player.Conn.Write([]byte{byte(MsgPing)}) // Simple pong
}

// handleJoinArena processes a join arena message
func handleJoinArena(msg *Message, player *Player, gs *GameState) {
	arenaID, team, err := ParseJoinArenaPacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse join arena: %v\n", err)
		return
	}

	arena := gs.ArenaManager.GetArena(arenaID)
	if arena == nil {
		fmt.Printf("Arena %d not found\n", arenaID)
		return
	}

	if arena.IsFull() {
		fmt.Printf("Arena %d is full\n", arenaID)
		return
	}

	err = arena.AddPlayer(player.ID, team)
	if err != nil {
		fmt.Printf("Failed to add player %d to arena %d: %v\n", player.ID, arenaID, err)
		return
	}

	fmt.Printf("Player %d joined arena %d as team %d\n", player.ID, arenaID, team)
}

// handleLeaveArena processes a leave arena message
func handleLeaveArena(msg *Message, player *Player, gs *GameState) {
	arenaID, err := ParseLeaveArenaPacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse leave arena: %v\n", err)
		return
	}

	arena := gs.ArenaManager.GetArena(arenaID)
	if arena == nil {
		fmt.Printf("Arena %d not found\n", arenaID)
		return
	}

	arena.RemovePlayer(player.ID)
	fmt.Printf("Player %d left arena %d\n", player.ID, arenaID)
}

// handleArenaList processes an arena list request
func handleArenaList(msg *Message, player *Player, gs *GameState) {
	gs.ArenaManager.mu.RLock()
	arenas := make([]*Arena, 0, len(gs.ArenaManager.Arenas))
	for _, arena := range gs.ArenaManager.Arenas {
		arenas = append(arenas, arena)
	}
	gs.ArenaManager.mu.RUnlock()

	// Send arena list to player
	response := fmt.Sprintf("Available arenas:\n")
	for _, arena := range arenas {
		response += fmt.Sprintf("- %s (ID: %d, Players: %d/%d)\n",
			arena.Name, arena.ID, arena.GetPlayerCount(), arena.MaxPlayers)
	}
	player.Conn.Write([]byte(response))
}

// handleArenaUpdate processes an arena update message
func handleArenaUpdate(msg *Message, player *Player, gs *GameState) {
	arenaID, x, y, err := ParseArenaUpdatePacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse arena update: %v\n", err)
		return
	}

	arena := gs.ArenaManager.GetArena(arenaID)
	if arena == nil {
		fmt.Printf("Arena %d not found\n", arenaID)
		return
	}

	arena.UpdatePlayerPosition(player.ID, x, y)
	fmt.Printf("Player %d updated position in arena %d: (%.2f, %.2f)\n", player.ID, arenaID, x, y)
}

// handleCastSpell processes a spell casting message
func handleCastSpell(msg *Message, player *Player, gs *GameState) {
	spellID, targetX, targetY, targetID, err := ParseCastSpellPacket(msg.Data)
	if err != nil {
		fmt.Printf("Failed to parse cast spell: %v\n", err)
		return
	}

	spellInstance, err := gs.SpellSystem.CastSpell(player.ID, spellID, targetX, targetY, targetID)
	if err != nil {
		fmt.Printf("Failed to cast spell %d: %v\n", spellID, err)
		return
	}

	fmt.Printf("Player %d cast spell %d at (%.2f, %.2f)\n", player.ID, spellID, targetX, targetY)
	_ = spellInstance // TODO: Broadcast spell cast to other players
}

// handleSpellList processes a spell list request
func handleSpellList(msg *Message, player *Player, gs *GameState) {
	spells := gs.SpellSystem.SpellManager.GetAllSpells()

	// Send spell list to player
	response := "Available spells:\n"
	for _, spell := range spells {
		response += fmt.Sprintf("- %s (ID: %d): %s\n", spell.Name, spell.ID, spell.Description)
	}
	player.Conn.Write([]byte(response))
}

// handleUnknownMessage captures unhandled packets for debugging
// handleUnknownMessage captures unhandled packets for debugging
func handleUnknownMessage(msg *Message, player *Player) {
	if debugLogger.Enabled {
		debugLogger.CapturePacket(msg, player)
	}
	fmt.Printf("Unknown message type: %d from player %s (ID: %d)\n", msg.Type, player.Name, player.ID)
}

// CapturePacket adds a packet to the debug log
func (d *DebugPacketLogger) CapturePacket(msg *Message, player *Player) {
	if !d.Enabled {
		return
	}

	capture := DebugPacketCapture{
		MessageType: msg.Type,
		Data:        make([]byte, len(msg.Data)),
		PlayerID:    player.ID,
		PlayerName:  player.Name,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
	}
	copy(capture.Data, msg.Data)

	d.CapturedPackets = append(d.CapturedPackets, capture)

	// Keep only the most recent packets
	if len(d.CapturedPackets) > d.MaxPackets {
		d.CapturedPackets = d.CapturedPackets[1:]
	}
}

// GetCapturedPackets returns all captured packets
func (d *DebugPacketLogger) GetCapturedPackets() []DebugPacketCapture {
	return d.CapturedPackets
}

// GetCapturedPacketsByType returns packets of a specific type
func (d *DebugPacketLogger) GetCapturedPacketsByType(msgType MessageType) []DebugPacketCapture {
	var filtered []DebugPacketCapture
	for _, packet := range d.CapturedPackets {
		if packet.MessageType == msgType {
			filtered = append(filtered, packet)
		}
	}
	return filtered
}

// ClearCapturedPackets clears all captured packets
func (d *DebugPacketLogger) ClearCapturedPackets() {
	d.CapturedPackets = make([]DebugPacketCapture, 0)
}

// GetDebugStats returns statistics about captured packets
func (d *DebugPacketLogger) GetDebugStats() map[MessageType]int {
	stats := make(map[MessageType]int)
	for _, packet := range d.CapturedPackets {
		stats[packet.MessageType]++
	}
	return stats
}

// sendCapturedPackets sends captured packets to the player
func sendCapturedPackets(player *Player) {
	packets := debugLogger.GetCapturedPackets()
	if len(packets) == 0 {
		player.Conn.Write([]byte("No captured packets\n"))
		return
	}

	response := fmt.Sprintf("Captured Packets (%d total):\n", len(packets))
	for i, packet := range packets {
		response += fmt.Sprintf("[%d] Type: %d, Player: %s (ID: %d), Time: %s, Data: %x\n",
			i+1, packet.MessageType, packet.PlayerName, packet.PlayerID, packet.Timestamp, packet.Data)
	}
	player.Conn.Write([]byte(response))
}

// sendDebugStats sends debug statistics to the player
func sendDebugStats(player *Player) {
	stats := debugLogger.GetDebugStats()
	if len(stats) == 0 {
		player.Conn.Write([]byte("No debug statistics available\n"))
		return
	}

	response := "Debug Statistics:\n"
	for msgType, count := range stats {
		response += fmt.Sprintf("Message Type %d: %d occurrences\n", msgType, count)
	}
	player.Conn.Write([]byte(response))
}

// sendDebugHelp sends help information for debug commands
func sendDebugHelp(player *Player) {
	help := `Debug Commands:
/debug packets - Show all captured unhandled packets
/debug stats - Show statistics of captured packets by type
/debug clear - Clear all captured packets
/debug help - Show this help message

Example: /debug packets
`
	player.Conn.Write([]byte(help))
}
