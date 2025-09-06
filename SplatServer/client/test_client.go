package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
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
)

// Team represents a team in the arena
type Team int

const (
	TeamNone   Team = 0
	TeamChaos  Team = 1
	TeamBalance Team = 2
	TeamOrder  Team = 3
)

func main() {
	fmt.Println("Arena Test Client: Connecting to SplatServer...")

	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected! Testing arena functionality...")

	// Send login message
	sendMessage(conn, MsgLogin, []byte("ArenaTester"))
	time.Sleep(1 * time.Second)

	// Request spell list
	sendMessage(conn, MsgSpellList, []byte{})
	time.Sleep(1 * time.Second)

	// Cast a spell (Fire Bolt at position 50, 50 targeting player 2)
	castData := make([]byte, 24)
	binary.LittleEndian.PutUint32(castData[0:4], 1) // spell ID (Fire Bolt)
	binary.LittleEndian.PutUint64(castData[4:12], uint64(50.0)) // target X
	binary.LittleEndian.PutUint64(castData[12:20], uint64(50.0)) // target Y
	binary.LittleEndian.PutUint32(castData[20:24], 2) // target player ID
	sendMessage(conn, MsgCastSpell, castData)
	time.Sleep(1 * time.Second)

	// Try to cast the same spell again (should be on cooldown)
	sendMessage(conn, MsgCastSpell, castData)
	time.Sleep(1 * time.Second)

	// Join arena 1 as Chaos team
	joinData := make([]byte, 8)
	binary.LittleEndian.PutUint32(joinData[0:4], 1) // arena ID
	binary.LittleEndian.PutUint32(joinData[4:8], uint32(TeamChaos)) // team
	sendMessage(conn, MsgJoinArena, joinData)
	time.Sleep(1 * time.Second)

	// Send arena position updates
	for i := 0; i < 5; i++ {
		updateData := make([]byte, 20)
		binary.LittleEndian.PutUint32(updateData[0:4], 1) // arena ID
		x := float64(10 + i)
		y := float64(20 + i*2)
		binary.LittleEndian.PutUint64(updateData[4:12], uint64(x))
		binary.LittleEndian.PutUint64(updateData[12:20], uint64(y))
		sendMessage(conn, MsgArenaUpdate, updateData)
		time.Sleep(500 * time.Millisecond)
	}

	// Leave arena
	leaveData := make([]byte, 4)
	binary.LittleEndian.PutUint32(leaveData[0:4], 1) // arena ID
	sendMessage(conn, MsgLeaveArena, leaveData)
	time.Sleep(1 * time.Second)

	// Send logout
	sendMessage(conn, MsgLogout, []byte{})

	fmt.Println("Arena test complete.")
}

func sendMessage(conn net.Conn, msgType MessageType, data []byte) {
	var buf bytes.Buffer

	// Write message type
	buf.WriteByte(byte(msgType))

	// Write data length
	length := uint32(len(data))
	binary.LittleEndian.PutUint32(make([]byte, 4), length)
	buf.Write(binary.LittleEndian.AppendUint32([]byte{}, length))

	// Write data
	buf.Write(data)

	conn.Write(buf.Bytes())
	fmt.Printf("Sent message type %d with %d bytes\n", msgType, len(data))
}
