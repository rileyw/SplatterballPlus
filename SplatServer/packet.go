package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// PacketType defines the type of packet
type PacketType uint16

const (
	PacketLogin         PacketType = 1
	PacketLogout        PacketType = 2
	PacketMove          PacketType = 3
	PacketChat          PacketType = 4
	PacketPing          PacketType = 5
	PacketPlayerUpdate  PacketType = 6
	PacketArenaUpdate   PacketType = 7
	PacketProjectile    PacketType = 8
	PacketSpellCast     PacketType = 9
	PacketGameState     PacketType = 10
)

// Packet represents a network packet
type Packet struct {
	Type   PacketType
	Length uint32
	Data   []byte
}

// NewPacket creates a new packet
func NewPacket(packetType PacketType, data []byte) *Packet {
	return &Packet{
		Type:   packetType,
		Length: uint32(len(data)),
		Data:   data,
	}
}

// Serialize converts the packet to bytes
func (p *Packet) Serialize() []byte {
	buf := new(bytes.Buffer)

	// Write packet type (2 bytes)
	binary.Write(buf, binary.LittleEndian, p.Type)

	// Write length (4 bytes)
	binary.Write(buf, binary.LittleEndian, p.Length)

	// Write data
	buf.Write(p.Data)

	return buf.Bytes()
}

// DeserializePacket reads a packet from a reader
func DeserializePacket(r io.Reader) (*Packet, error) {
	var packetType PacketType
	if err := binary.Read(r, binary.LittleEndian, &packetType); err != nil {
		return nil, err
	}

	var length uint32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	return &Packet{
		Type:   packetType,
		Length: length,
		Data:   data,
	}, nil
}

// Packet builders for common packets
func BuildLoginPacket(name string) *Packet {
	data := []byte(name)
	return NewPacket(PacketLogin, data)
}

func BuildMovePacket(x, y float64) *Packet {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, x)
	binary.Write(buf, binary.LittleEndian, y)
	return NewPacket(PacketMove, buf.Bytes())
}

func BuildChatPacket(message string) *Packet {
	data := []byte(message)
	return NewPacket(PacketChat, data)
}

func BuildPlayerUpdatePacket(player *Player) *Packet {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, player.ID)
	buf.WriteString(player.Name)
	buf.WriteByte(0) // null terminator
	binary.Write(buf, binary.LittleEndian, player.X)
	binary.Write(buf, binary.LittleEndian, player.Y)
	binary.Write(buf, binary.LittleEndian, player.Health)
	return NewPacket(PacketPlayerUpdate, buf.Bytes())
}

// Packet parsers
func ParseLoginPacket(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty login data")
	}
	return string(data), nil
}

func ParseMovePacket(data []byte) (float64, float64, error) {
	if len(data) < 16 {
		return 0, 0, fmt.Errorf("insufficient move data")
	}
	var x, y float64
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &x)
	binary.Read(buf, binary.LittleEndian, &y)
	return x, y, nil
}

func ParseChatPacket(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty chat data")
	}
	return string(data), nil
}

func ParseJoinArenaPacket(data []byte) (int, Team, error) {
	if len(data) < 8 {
		return 0, TeamNone, fmt.Errorf("insufficient join arena data")
	}
	var arenaID int
	var team Team
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &arenaID)
	binary.Read(buf, binary.LittleEndian, &team)
	return arenaID, team, nil
}

func ParseLeaveArenaPacket(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient leave arena data")
	}
	var arenaID int
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &arenaID)
	return arenaID, nil
}

func ParseArenaUpdatePacket(data []byte) (int, float64, float64, error) {
	if len(data) < 20 {
		return 0, 0, 0, fmt.Errorf("insufficient arena update data")
	}
	var arenaID int
	var x, y float64
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &arenaID)
	binary.Read(buf, binary.LittleEndian, &x)
	binary.Read(buf, binary.LittleEndian, &y)
	return arenaID, x, y, nil
}

func ParseCastSpellPacket(data []byte) (int, float64, float64, int, error) {
	if len(data) < 24 {
		return 0, 0, 0, 0, fmt.Errorf("insufficient cast spell data")
	}
	var spellID int
	var targetX, targetY float64
	var targetID int
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &spellID)
	binary.Read(buf, binary.LittleEndian, &targetX)
	binary.Read(buf, binary.LittleEndian, &targetY)
	binary.Read(buf, binary.LittleEndian, &targetID)
	return spellID, targetX, targetY, targetID, nil
}
