package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

const (
	SERVER_PORT = "4000"
)

var (
	playerIDCounter int64
	gameState       *GameState
)

func main() {
	fmt.Println("SplatServer: Starting game server on port", SERVER_PORT)

	// Initialize database
	InitDB()
	CreateTables()

	gameState = NewGameState()

	// Initialize spell system
	InitializeSpellSystem()

	// Initialize basic arenas
	initializeArenas(gameState)

	// Start the game loop in a goroutine
	go GameLoop(gameState)

	// Start UDP listener for real-time updates
	go startUDPListener()

	ln, err := net.Listen("tcp", ":"+SERVER_PORT)
	if err != nil {
		log.Fatalf("Error starting TCP server: %v", err)
	}
	defer ln.Close()
	fmt.Println("SplatServer: TCP server listening for connections...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func startUDPListener() {
	addr, err := net.ResolveUDPAddr("udp", ":"+SERVER_PORT)
	if err != nil {
		log.Printf("Failed to resolve UDP address: %v", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Printf("Failed to start UDP listener: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("SplatServer: UDP listener started")

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("UDP read error: %v", err)
			continue
		}

		// Handle UDP packet
		go handleUDPPacket(conn, addr, buf[:n])
	}
}

func handleUDPPacket(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	// For now, just echo back
	conn.WriteToUDP(data, addr)
	fmt.Printf("UDP packet from %s: %x\n", addr.String(), data)
}

func initializeArenas(gs *GameState) {
	// Create some basic arenas
	gs.ArenaManager.CreateArena(1, "Chaos Arena", 8, 1)
	gs.ArenaManager.CreateArena(2, "Balance Arena", 8, 2)
	gs.ArenaManager.CreateArena(3, "Order Arena", 8, 3)

	fmt.Println("SplatServer: Initialized 3 basic arenas")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	playerID := int(atomic.AddInt64(&playerIDCounter, 1))

	player := &Player{
		ID:       playerID,
		Name:     fmt.Sprintf("Player%d", playerID),
		X:        0,
		Y:        0,
		Health:   100,
		Conn:     conn,
		LastSeen: time.Now(),
	}

	gameState.AddPlayer(player)
	fmt.Printf("New connection from %s (Player ID: %d)\n", clientAddr, playerID)
	conn.Write([]byte(fmt.Sprintf("Welcome %s! Your ID is %d\n", player.Name, playerID)))

	// Handle incoming messages
	for {
		msg, err := ParseMessage(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error parsing message from %d: %v\n", playerID, err)
			}
			gameState.RemovePlayer(playerID)
			return
		}
		HandleMessage(msg, player, gameState)
		player.LastSeen = time.Now()
	}
}
