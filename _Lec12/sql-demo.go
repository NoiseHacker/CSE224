package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 1. Open (or create) the database file
	db, err := sql.Open("sqlite3", "demo.db")
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// 2. Create a table if it doesn’t exist
	createStmt := `
    CREATE TABLE IF NOT EXISTS users (
      id   INTEGER PRIMARY KEY,
      name TEXT NOT NULL,
      age  INTEGER
    );`
	if _, err := db.Exec(createStmt); err != nil {
		log.Fatalf("create table: %v", err)
	}

	// 3. Insert two rows
	if _, err := db.Exec(
		"INSERT INTO users(name, age) VALUES (?, ?)",
		"Alice", 30,
	); err != nil {
		log.Fatalf("insert: %v", err)
	}

	if _, err := db.Exec(
		"INSERT INTO users(name, age) VALUES (?, ?)",
		"Bob", 45,
	); err != nil {
		log.Fatalf("insert: %v", err)
	}

	// 4. Query rows
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatalf("query: %v", err)
	}
	defer rows.Close()

	// 5. Iterate and scan
	for rows.Next() {
		var id, age int
		var name string
		if err := rows.Scan(&id, &name, &age); err != nil {
			log.Fatalf("scan: %v", err)
		}
		log.Printf("User #%d: %s (age %d)", id, name, age)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("rows: %v", err)
	}
}
