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

	// Création de la table votes
	votesStatement := `
	CREATE TABLE IF NOT EXISTS votes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		username TEXT,
		vote_type INTEGER,
		FOREIGN KEY (post_id) REFERENCES posts(id),
		UNIQUE(post_id, username)
	);`
	_, err = db.Exec(votesStatement)
	if err != nil {
		panic(err)
	}
}
