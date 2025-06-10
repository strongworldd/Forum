package main

import (
	"database/sql"
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func main() {
	database, _ := sql.Open("sqlite3", "./foo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (name) VALUES (?)")
	statement.Exec("John Doe")
	rows, _ := database.Query("SELECT id, name FROM people")
	var id int 
	var name string
	for rows.Next() {
		rows.Scan(&id, &name)
		fmt.Println(strconv.Itoa(id) + ": " + name)
	}
}