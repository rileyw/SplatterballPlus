# SplatServer (Go Version)

A minimal, cross-platform game server skeleton written in Go. This project is the starting point for migrating MageServer to Go for efficient cloud hosting.

## Features
- TCP server listening on port 4000
- UDP server for real-time updates on port 4000
- Accepts multiple client connections
- Greets clients on connect
- Foundation for custom game protocol and logic
- **Entity structs**: Player, GameState, Arena with thread-safe management
- **Game loop**: Runs at 60 ticks/second, updates player health, checks timeouts
- **Protocol handling**: Parses binary messages (login, move, chat, logout, ping)
- **Advanced packet system**: Structured packets with types, serialization, and parsing
- **Persistence**: MySQL database for player data (auto-saves on login/logout/timeout)
- **In-memory database**: SQLite support for development and testing

## Prerequisites
- Go 1.20 or newer (https://golang.org/dl/)
- MySQL 5.7+ or compatible database (for production), or SQLite (for development)

## Database Setup
1. Install MySQL (for production) or use SQLite (for development).
2. For MySQL: Create a database:
   ```sql
   CREATE DATABASE splatserver;
   ```
3. Set environment variables:
   ```sh
   # For MySQL (default)
   export DB_TYPE=mysql
   export DB_USER=your_user
   export DB_PASSWORD=your_password
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_NAME=splatserver

   # For in-memory SQLite (development)
   export DB_TYPE=sqlite
   ```
   For cloud databases (e.g., AWS RDS, Google Cloud SQL), update DB_HOST accordingly.

The server will automatically create the necessary tables on startup.

## Build Instructions

1. Clone or copy the SplatServer directory:
   ```sh
   cd SplatServer
   ```
2. Build the server binary:
   ```sh
   go build -o splatserver main.go
   ```
3. Run the server:
   ```sh
   ./splatserver
   ```
   The server will start on TCP port 4000.

## Cloud Deployment Tips
- Use a small cloud VM (e.g., AWS Lightsail, DigitalOcean, GCP, Azure) with Go installed.
- Open port 4000 in your firewall/security group.
- For production, run the server as a systemd service or in a Docker container.

## Testing

### Full Test Suite
Run all tests (unit + integration) with the provided script:
```sh
./run_tests.sh
```
This will run unit tests, start the server, run the test client, and stop the server.

### Unit Tests Only
Run unit tests for entities, game loop, and protocol:
```sh
go test -v
```

### Integration Testing Only
1. Start the server in one terminal:
   ```sh
   go run .
   ```
2. Run the test client in another terminal:
   ```sh
   cd client && go run test_client.go
   ```
   This will connect and send sample messages (login, move, chat, ping, logout).

### Manual Testing
Use tools like `netcat` or `telnet` to send raw binary data, or write a custom client.

## Debug Capabilities

The server includes built-in debug capabilities for capturing and analyzing unhandled packets:

### Configuration
Debug packet capture is controlled by environment variables:

```sh
# Enable debug packet capture (default: disabled)
export DEBUG_PACKETS_ENABLED=true

# Maximum number of packets to store (default: 100)
export DEBUG_PACKETS_MAX=50
```

**Note**: Debug packet capture is disabled by default for performance reasons. Only enable in development or debugging scenarios.

**Security Note**: Captured packets may contain sensitive data. Ensure debug features are only enabled in secure development environments.

### Debug Commands
Send these commands as chat messages to access debug features:

- `/debug packets` - Display all captured unhandled packets with details
- `/debug stats` - Show statistics of captured packets by message type
- `/debug clear` - Clear all captured packets from memory
- `/debug help` - Show available debug commands

### Packet Capture
- Automatically captures incoming packets with unknown message types
- Stores up to the configured maximum number of recent unhandled packets
- Records timestamp, player info, message type, and raw data
- Useful for reverse-engineering client protocols or debugging new features

### Example Output
```
Captured Packets (2 total):
[1] Type: 15, Player: TestPlayer (ID: 1), Time: 2025-09-05 14:30:15, Data: 0102030405
[2] Type: 20, Player: TestPlayer (ID: 1), Time: 2025-09-05 14:30:20, Data: 0a0b0c0d0e
```

## License
MIT (or match your main project license)
