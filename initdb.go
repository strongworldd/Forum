package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	initAccountsDB()
	initPostsDB()
	fmt.Println("Bases de données initialisées avec succès ✅")
}

func initAccountsDB() {
	db, err := sql.Open("sqlite3", "./BDD/accounts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	statement := `
	CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		password TEXT,
		email TEXT,
		postliked TEXT,
		uuid TEXT
	);`
	_, err = db.Exec(statement)
	if err != nil {
		panic(err)
	}
}

func initPostsDB() {
	db, err := sql.Open("sqlite3", "./BDD/posts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	statement := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		author TEXT
	);`
	_, err = db.Exec(statement)
	if err != nil {
		panic(err)
	}
}