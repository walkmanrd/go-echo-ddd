package configs

import (
	"database/sql"
	"log"
	"os"
)

// DB is an interface for database
type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
}

// sqlOpener is a function that open database
type (
	sqlOpener func(string, string) (*sql.DB, error)
)

// AutoMigrate is a function that create table if not exist
func AutoMigrate(db DB) {
	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err := db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table: ", err)
	}
}

// OpenDB is a function that open database
func OpenDB(open sqlOpener, connectionUrl string) (*sql.DB, error) {
	return open("postgres", connectionUrl)
}

// ConnectDatabase is a function that connect database
func ConnectDatabase() *sql.DB {
	var err error

	db, err := OpenDB(sql.Open, os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Connect database error", err)
	}

	return db
}
