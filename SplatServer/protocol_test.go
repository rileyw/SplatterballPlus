package main

import (
	"net"
	"testing"
)

func TestDebugPacketCapture(t *testing.T) {
	// Enable debug logging for this test
	originalEnabled := debugLogger.Enabled
	debugLogger.Enabled = true
	defer func() { debugLogger.Enabled = originalEnabled }()

	// Create a test player
	player := &Player{
		ID:   1,
		Name: "TestPlayer",
		Conn: &net.TCPConn{}, // Mock connection
	}

	// Create a test message
	msg := &Message{
		Type: 99, // Unknown message type
		Data: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	}

	// Capture the packet
	debugLogger.CapturePacket(msg, player)

	// Verify the packet was captured
	packets := debugLogger.GetCapturedPackets()
	if len(packets) != 1 {
		t.Errorf("Expected 1 captured packet, got %d", len(packets))
	}

	captured := packets[0]
	if captured.MessageType != 99 {
		t.Errorf("Expected message type 99, got %d", captured.MessageType)
	}
	if captured.PlayerID != 1 {
		t.Errorf("Expected player ID 1, got %d", captured.PlayerID)
	}
	if captured.PlayerName != "TestPlayer" {
		t.Errorf("Expected player name 'TestPlayer', got '%s'", captured.PlayerName)
	}
	if len(captured.Data) != 5 {
		t.Errorf("Expected data length 5, got %d", len(captured.Data))
	}

	// Test stats
	stats := debugLogger.GetDebugStats()
	if stats[99] != 1 {
		t.Errorf("Expected 1 occurrence of message type 99, got %d", stats[99])
	}

	// Test filtering by type
	filtered := debugLogger.GetCapturedPacketsByType(99)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered packet, got %d", len(filtered))
	}

	// Test clear
	debugLogger.ClearCapturedPackets()
	packets = debugLogger.GetCapturedPackets()
	if len(packets) != 0 {
		t.Errorf("Expected 0 packets after clear, got %d", len(packets))
	}
}

func TestDebugPacketCaptureDisabled(t *testing.T) {
	// Disable debug logging for this test
	originalEnabled := debugLogger.Enabled
	debugLogger.Enabled = false
	defer func() { debugLogger.Enabled = originalEnabled }()

	// Clear any existing packets
	debugLogger.ClearCapturedPackets()

	// Create a test player
	player := &Player{
		ID:   1,
		Name: "TestPlayer",
		Conn: &net.TCPConn{}, // Mock connection
	}

	// Create a test message
	msg := &Message{
		Type: 99, // Unknown message type
		Data: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	}

	// Try to capture the packet
	debugLogger.CapturePacket(msg, player)

	// Verify no packet was captured
	packets := debugLogger.GetCapturedPackets()
	if len(packets) != 0 {
		t.Errorf("Expected 0 captured packets when disabled, got %d", len(packets))
	}
}

func TestDebugPacketMaxLimit(t *testing.T) {
	// Enable debug logging for this test
	originalEnabled := debugLogger.Enabled
	debugLogger.Enabled = true
	defer func() { debugLogger.Enabled = originalEnabled }()

	debugLogger.ClearCapturedPackets()
	debugLogger.MaxPackets = 3

	player := &Player{ID: 1, Name: "TestPlayer"}

	// Add 5 packets (more than the limit)
	for i := 0; i < 5; i++ {
		msg := &Message{
			Type: MessageType(i + 100),
			Data: []byte{byte(i)},
		}
		debugLogger.CapturePacket(msg, player)
	}

	packets := debugLogger.GetCapturedPackets()
	if len(packets) != 3 {
		t.Errorf("Expected 3 packets (max limit), got %d", len(packets))
	}

	// Verify we kept the most recent packets (types 102, 103, 104)
	expectedTypes := []MessageType{102, 103, 104}
	for i, packet := range packets {
		if packet.MessageType != expectedTypes[i] {
			t.Errorf("Expected packet type %d at position %d, got %d", expectedTypes[i], i, packet.MessageType)
		}
	}
}
