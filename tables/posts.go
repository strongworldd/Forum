package tables

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
    "os"
    "net/http"
    "encoding/json"
)

func ResetPostsTable() {
	database, _ := sql.Open("sqlite3", "../BDD/posts.db")
    defer database.Close()
    _, _ = database.Exec("DELETE FROM posts")
}

func Deletepost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
        return
    }
    var data struct {
        ID int `json:"id"`
    }
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Erreur de décodage JSON", http.StatusBadRequest)
        return
    }
    id := data.ID
	database, err := sql.Open("sqlite3", "../BDD/posts.db")
    if err != nil {
        fmt.Println("Erreur ouverture BDD:", err)
        return
    }
    defer database.Close()

    rows, err := database.Query("SELECT content FROM posts WHERE id = ?", id)

    if err != nil {
        fmt.Println("Erreur récupération images post:", err)
    } else {
        var image string
        for rows.Next() {
            err = rows.Scan(&image)
            if err != nil {
                fmt.Println("Erreur scan image:", err)
                return
            }
        }

        err := os.Remove("../img/" + image)
        if err != nil {
            fmt.Println("Erreur suppression image:", err)
        }
        rows.Close()
    }

    _, err = database.Exec("DELETE FROM posts WHERE id = ?", id)
    if err != nil {
        fmt.Println("Erreur suppression post:", err)
    }
}

func CreatePost() {
	database, _ := sql.Open("sqlite3", "../BDD/posts.db")
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT)")
	statement.Exec()

	statement, _ = database.Prepare("INSERT INTO posts (title, content, author) VALUES (?, ?, ?)")
	statement.Exec("First Post", "This is the content of the first post.", "John Doe")
}

func CheckPostDB() {
	database, _ := sql.Open("sqlite3", "../BDD/posts.db")
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT)")
	statement.Exec()
}

func LoadPosts() {
	database, _ := sql.Open("sqlite3", "../BDD/posts.db")
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