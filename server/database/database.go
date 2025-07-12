package database

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

type database struct {
	db *sql.DB
}

func (d *database) GetDB() *sql.DB {
	if d.db == nil {
		d.db = initDatabase()
	}
	return d.db
}

var DB *database = &database{}

func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite", "./database/exchanges.db")
	if err != nil {
		slog.Error("Error to open `exchanges.db` database")
		panic(err)
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS exchanges (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		rate_value TEXT
	)`); err != nil {
		slog.Error("Erro to create database table `exchanges`")
		panic(err)
	}
	slog.Info("Database `exchanges.db` ready")
	return db
}
