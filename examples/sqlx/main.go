package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/mtsmfm/sqlcodegen/examples/sqlx/sqlstructs"
)

//go:generate go run github.com/mtsmfm/sqlcodegen g
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

	db.Exec(`
		INSERT INTO users VALUES (1, 'a')
	`)

	var results1 []sqlstructs.X
	// sqlcodegen X
	err = db.Select(&results1, "SELECT id FROM users LIMIT 5")
	if err != nil {
		panic(err)
	}
	for _, result := range results1 {
		log.Printf("%+v", result)
	}

	var results2 []sqlstructs.Foo
	// sqlcodegen Foo
	err = db.Select(&results2, "SELECT firebase_auth_uid FROM users LIMIT 5")
	if err != nil {
		panic(err)
	}
	for _, result := range results2 {
		log.Printf("%+v", result)
	}

	var results3 []sqlstructs.Bar
	// sqlcodegen Bar
	err = db.Select(&results3, `
		SELECT *
		FROM
		users
		LIMIT 5
	`)

	if err != nil {
		panic(err)
	}
	for _, result := range results3 {
		log.Printf("%+v", result)
	}
}
