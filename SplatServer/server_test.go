package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"
)

// TestGameState tests player management in GameState
func TestGameState(t *testing.T) {
	gs := NewGameState()

	player := &Player{ID: 1, Name: "TestPlayer"}
	gs.AddPlayer(player)

	if _, exists := gs.GetPlayer(1); !exists {
		t.Error("Player should exist after adding")
	}

	gs.RemovePlayer(1)
	if _, exists := gs.GetPlayer(1); exists {
		t.Error("Player should not exist after removing")
	}
}

// TestParseMessage tests message parsing
func TestParseMessage(t *testing.T) {
	// Simulate a login message using new packet format
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, PacketType(MsgLogin)) // 2 bytes
	binary.Write(&buf, binary.LittleEndian, uint32(5))            // 4 bytes length
	buf.WriteString("Alice")                                       // 5 bytes data

	conn := &mockConn{data: buf.Bytes()}
	msg, err := ParseMessage(conn)
	if err != nil {
		t.Fatalf("ParseMessage failed: %v", err)
	}
	if msg.Type != MsgLogin {
		t.Errorf("Expected MsgLogin, got %v", msg.Type)
	}
	if string(msg.Data) != "Alice" {
		t.Errorf("Expected 'Alice', got %s", string(msg.Data))
	}
}

// mockConn simulates a net.Conn for testing
type mockConn struct {
	data []byte
	pos  int
}

func (m *mockConn) Read(b []byte) (int, error) {
	if m.pos >= len(m.data) {
		return 0, io.EOF
	}
	n := copy(b, m.data[m.pos:])
	m.pos += n
	return n, nil
}

func (m *mockConn) Write(b []byte) (int, error) { return len(b), nil }
func (m *mockConn) Close() error               { return nil }
func (m *mockConn) LocalAddr() net.Addr         { return nil }
func (m *mockConn) RemoteAddr() net.Addr        { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

// TestGameLoop tests the game loop updates
func TestGameLoop(t *testing.T) {
	gs := NewGameState()
	player := &Player{ID: 1, Name: "Test", Health: 50, LastSeen: time.Now()}
	gs.AddPlayer(player)

	// Run update manually
	UpdateGameState(gs)

	if player.Health != 51 {
		t.Errorf("Expected health 51, got %d", player.Health)
	}
}

// Integration test: Start server and connect a client
func TestServerIntegration(t *testing.T) {
	// This test requires the server to be running
	// For simplicity, skip if not running
	t.Skip("Integration test requires running server")

	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send a login message
	var buf bytes.Buffer
	buf.WriteByte(byte(MsgLogin))
	binary.Write(&buf, binary.LittleEndian, uint32(4))
	buf.WriteString("Bob")

	conn.Write(buf.Bytes())

	// Read response (if any)
	response := make([]byte, 1024)
	conn.Read(response)
	// TODO: Assert response
}
