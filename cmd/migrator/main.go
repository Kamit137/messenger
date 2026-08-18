package main

import (
	"flag"
	"log"
	"messenger/internal/storage"
	"os"
	"path/filepath"
)

func main() {
	databasePath := flag.String("database-path", "storage/messenger.db", "path to SQLite database")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*databasePath), 0o755); err != nil {
		log.Fatalf("create database directory: %v", err)
	}
	err := storage.NewSQLiteStorage(*databasePath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer storage.Close()
	if err := storage.Migrate(); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	log.Println("database schema is ready")
}
