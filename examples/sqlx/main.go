package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id bigint PRIMARY KEY,
			firebase_auth_uid text NOT NULL UNIQUE
		);
	`)

	if err != nil {
		panic(err)
	}

	var results []interface{}
	err = db.Select(&results, "SELECT id FROM users LIMIT 5")

	if err != nil {
		panic(err)
	}

	for _, result := range results {
		log.Printf("%v", result)
	}
}
