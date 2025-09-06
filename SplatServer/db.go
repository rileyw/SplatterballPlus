package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var dbType string

// InitDB initializes the database connection
func InitDB() {
	dbType = getEnv("DB_TYPE", "mysql") // Default to MySQL, or "sqlite" for in-memory

	switch dbType {
	case "sqlite":
		initSQLite()
	case "mysql":
		initMySQL()
	default:
		log.Fatalf("Unsupported DB_TYPE: %s", dbType)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Printf("Connected to %s database\n", dbType)
}

// initSQLite initializes an in-memory SQLite database
func initSQLite() {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
}

// initMySQL initializes a MySQL database connection
func initMySQL() {
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "splatserver")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}
}

// CreateTables creates necessary tables if they don't exist
func CreateTables() {
	var playerTable string
	if dbType == "sqlite" {
		playerTable = `
		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			x REAL DEFAULT 0,
			y REAL DEFAULT 0,
			health INTEGER DEFAULT 100,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		playerTable = `
		CREATE TABLE IF NOT EXISTS players (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			x FLOAT DEFAULT 0,
			y FLOAT DEFAULT 0,
			health INT DEFAULT 100,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`
	}

	_, err := db.Exec(playerTable)
	if err != nil {
		log.Fatalf("Failed to create players table: %v", err)
	}

	fmt.Println("Database tables ready")
}

// SavePlayer saves or updates a player in the database
func SavePlayer(p *Player) error {
	if dbType == "sqlite" {
		// SQLite doesn't support ON DUPLICATE KEY UPDATE
		// First try to update, if no rows affected, insert
		result, err := db.Exec(`
			UPDATE players SET name=?, x=?, y=?, health=?, updated_at=datetime('now')
			WHERE id=?`, p.Name, p.X, p.Y, p.Health, p.ID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			// No existing record, insert new one
			_, err = db.Exec(`
				INSERT INTO players (id, name, x, y, health, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
				p.ID, p.Name, p.X, p.Y, p.Health)
			return err
		}

		return nil
	} else {
		// MySQL version
		query := `
			INSERT INTO players (id, name, x, y, health)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE name=VALUES(name), x=VALUES(x), y=VALUES(y), health=VALUES(health)`

		_, err := db.Exec(query, p.ID, p.Name, p.X, p.Y, p.Health)
		return err
	}
}

// LoadPlayer loads a player from the database by ID
func LoadPlayer(id int) (*Player, error) {
	var p Player
	query := "SELECT id, name, x, y, health FROM players WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&p.ID, &p.Name, &p.X, &p.Y, &p.Health)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
