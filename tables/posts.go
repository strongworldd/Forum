package tables

import (
	"database/sql"
	"fmt"
	//"strconv"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func ResetPostsTable() {
	database, _ := sql.Open("sqlite3", "./BDD/posts.db")
    defer database.Close()
    _, _ = database.Exec("DELETE FROM posts")
}

func Deletepost(id int) {
	database, _ := sql.Open("sqlite3", "./BDD/posts.db")
    defer database.Close()
    _, _ = database.Exec("DELETE FROM posts WHERE id = ?", id)
}

func CreatePost() {
	database, _ := sql.Open("sqlite3", "./BDD/posts.db")
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT)")
	statement.Exec()

	statement, _ = database.Prepare("INSERT INTO posts (title, content, author) VALUES (?, ?, ?)")
	statement.Exec("First Post", "This is the content of the first post.", "John Doe")
}

func CheckPostDB() {
	database, _ := sql.Open("sqlite3", "./BDD/posts.db")
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT)")
	statement.Exec()
}

func LoadPosts() {
	database, _ := sql.Open("sqlite3", "./BDD/posts.db")
	defer database.Close()

	rows, _ := database.Query("SELECT id, title, content, author FROM posts")
	defer rows.Close()

	var id int
	var title string
	var content string
	var author string
	for rows.Next() {
		rows.Scan(&id, &title, &content, &author)
		fmt.Printf("%d: %s by %s\nContent: %s\n", id, title, author, content)
	}
}