package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

// TestPacketSerialization tests packet serialization and deserialization
func TestPacketSerialization(t *testing.T) {
	// Test login packet
	original := BuildLoginPacket("TestPlayer")
	serialized := original.Serialize()

	// Deserialize
	reader := bytes.NewReader(serialized)
	deserialized, err := DeserializePacket(reader)
	if err != nil {
		t.Fatalf("Failed to deserialize packet: %v", err)
	}

	if deserialized.Type != original.Type {
		t.Errorf("Type mismatch: expected %v, got %v", original.Type, deserialized.Type)
	}

	if deserialized.Length != original.Length {
		t.Errorf("Length mismatch: expected %v, got %v", original.Length, deserialized.Length)
	}

	if !bytes.Equal(deserialized.Data, original.Data) {
		t.Errorf("Data mismatch: expected %v, got %v", original.Data, deserialized.Data)
	}
}

// TestPacketBuilders tests all packet builder functions
func TestPacketBuilders(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() *Packet
		expected PacketType
	}{
		{"Login", func() *Packet { return BuildLoginPacket("test") }, PacketLogin},
		{"Move", func() *Packet { return BuildMovePacket(10.5, 20.3) }, PacketMove},
		{"Chat", func() *Packet { return BuildChatPacket("hello") }, PacketChat},
		{"PlayerUpdate", func() *Packet {
			player := &Player{ID: 1, Name: "Test", X: 1.0, Y: 2.0, Health: 100}
			return BuildPlayerUpdatePacket(player)
		}, PacketPlayerUpdate},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packet := tt.builder()
			if packet.Type != tt.expected {
				t.Errorf("Expected packet type %v, got %v", tt.expected, packet.Type)
			}
			if packet.Length == 0 {
				t.Error("Packet length should not be zero")
			}
		})
	}
}

// TestPacketParsers tests packet parsing functions
func TestPacketParsers(t *testing.T) {
	// Test login parser
	loginData := []byte("TestPlayer")
	name, err := ParseLoginPacket(loginData)
	if err != nil {
		t.Fatalf("ParseLoginPacket failed: %v", err)
	}
	if name != "TestPlayer" {
		t.Errorf("Expected 'TestPlayer', got '%s'", name)
	}

	// Test move parser
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, 10.5)
	binary.Write(&buf, binary.LittleEndian, 20.3)
	x, y, err := ParseMovePacket(buf.Bytes())
	if err != nil {
		t.Fatalf("ParseMovePacket failed: %v", err)
	}
	if x != 10.5 || y != 20.3 {
		t.Errorf("Expected (10.5, 20.3), got (%f, %f)", x, y)
	}

	// Test chat parser
	chatData := []byte("Hello world")
	message, err := ParseChatPacket(chatData)
	if err != nil {
		t.Fatalf("ParseChatPacket failed: %v", err)
	}
	if message != "Hello world" {
		t.Errorf("Expected 'Hello world', got '%s'", message)
	}
}

// TestPacketErrors tests error handling in packet operations
func TestPacketErrors(t *testing.T) {
	// Test parsing empty login packet
	_, err := ParseLoginPacket([]byte{})
	if err == nil {
		t.Error("Expected error for empty login data")
	}

	// Test parsing insufficient move data
	_, _, err = ParseMovePacket([]byte{1, 2, 3})
	if err == nil {
		t.Error("Expected error for insufficient move data")
	}

	// Test parsing empty chat packet
	_, err = ParseChatPacket([]byte{})
	if err == nil {
		t.Error("Expected error for empty chat data")
	}
}

// TestConcurrentPackets tests packet operations under concurrent access
func TestConcurrentPackets(t *testing.T) {
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			// Test serialization
			packet := BuildLoginPacket(fmt.Sprintf("Player%d", id))
			serialized := packet.Serialize()

			// Test deserialization
			reader := bytes.NewReader(serialized)
			deserialized, err := DeserializePacket(reader)
			if err != nil {
				t.Errorf("Concurrent deserialize failed: %v", err)
			}

			if deserialized.Type != PacketLogin {
				t.Errorf("Concurrent type check failed")
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
