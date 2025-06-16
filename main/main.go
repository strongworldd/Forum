package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/tables"
	"io"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "../BDD/accounts.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/createpost", createPostHandler)
	http.HandleFunc("/api/posts", postsAPIHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("../img"))))
	http.Handle("/Login/", http.StripPrefix("/Login/", http.FileServer(http.Dir("../Login"))))
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("../html"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/html/home.html", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	fmt.Println("Serveur démarré sur : http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	hashedPassword, err := tables.HashPassword(data.Password)
	if err != nil {
		http.Error(w, "Erreur de hash du mot de passe", 500)
		return
	}

	err = tables.NewUserRepository(db).CreateUser(data.Username, hashedPassword, data.Email)
	if err != nil {
		http.Error(w, "Erreur création utilisateur", 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	repo := tables.NewUserRepository(db)
	hash, err := repo.GetPasswordByUsername(data.Username)
	if err != nil || hash == "" {
		http.Error(w, "Utilisateur inconnu", http.StatusUnauthorized)
		return
	}
	if !tables.ComparePasswords(hash, data.Password) {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")

	file, handler, err := r.FormFile("content")
	if err != nil {
		http.Error(w, "Erreur fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imgPath := "../img/" + handler.Filename
	dst, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Erreur sauvegarde image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", 500)
		return
	}
	defer dbPosts.Close()

	stmt, err := dbPosts.Prepare("INSERT INTO posts (title, content, author) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Erreur préparation", 500)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, handler.Filename, author)
	if err != nil {
		http.Error(w, "Erreur insertion", 500)
		return
	}

	http.Redirect(w, r, "/html/home copy.html", http.StatusSeeOther)
}

func postsAPIHandler(w http.ResponseWriter, r *http.Request) {
	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", 500)
		return
	}
	defer dbPosts.Close()

	rows, err := dbPosts.Query("SELECT title, content, author FROM posts ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Erreur lecture BDD", 500)
		return
	}
	defer rows.Close()

	type Post struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.Title, &p.Content, &p.Author)
		posts = append(posts, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
