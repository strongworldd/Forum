package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/tables"
	"io"
	"net/http"
	"os"
	"strconv"

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
	http.HandleFunc("/api/posts", postsAPIHandler) // <-- Ajout de la route API pour les posts
	http.HandleFunc("/deletepost", tables.Deletepost)
	http.HandleFunc("/api/me", getuserdata)

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

	id, err := repo.GetUserID(data.Username)
	if err != nil || id == 0 {
		http.Error(w, "Erreur récupération ID utilisateur", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionid",
		Value:    strconv.Itoa(id), // par exemple un UUID
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func getuserdata(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("sessionid")
    if err != nil {
        http.Error(w, "Non authentifié", http.StatusUnauthorized)
        return
    }
    userID, err := strconv.Atoi(cookie.Value)
    if err != nil {
        http.Error(w, "Session invalide", http.StatusUnauthorized)
        return
    }
	
    user, err := tables.GetUserByID(userID)
    if err != nil {
        http.Error(w, "Utilisateur inconnu", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("sessionid")
    if err != nil {
        http.Error(w, "Non authentifié", http.StatusUnauthorized)
        return
    }
    author, _ := strconv.Atoi(cookie.Value)

	title := r.FormValue("title")

	// On récupère le fichier image
	file, _, err := r.FormFile("content")
	if err != nil {
		http.Error(w, "Erreur fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// On enregistre le chemin de l'image dans la BDD
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

	_, err = stmt.Exec(title, "post0.jpg", author)
	if err != nil {
		http.Error(w, "Erreur insertion", 500)
		return
	}

	http.Redirect(w, r, "/html/home copy.html", http.StatusSeeOther)

	rows, _ := dbPosts.Query("SELECT id FROM posts")

	var postID int
	for rows.Next() {
		err = rows.Scan(&postID)
		if err != nil {
			http.Error(w, "Erreur lecture ID post", http.StatusInternalServerError)
			return
		}
	}

	filename := "post" + strconv.Itoa(postID) + ".jpg"

	// On sauve le fichier dans le dossier img/
	imgPath := "../img/" + filename
	dst, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Erreur sauvegarde image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// On met à jour le chemin de l'image dans la BDD
	_, _ = dbPosts.Exec("UPDATE posts SET content = ? WHERE id = ?", filename, postID)
}

func postsAPIHandler(w http.ResponseWriter, r *http.Request) {
	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", 500)
		return
	}
	defer dbPosts.Close()

	rows, err := dbPosts.Query("SELECT id, title, content, author FROM posts ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Erreur lecture BDD", 500)
		return
	}
	defer rows.Close()

	type Post struct {
		ID 		int `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	var posts []Post
	for rows.Next() {
		var p Post
		rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author)
		id, _ := strconv.Atoi(p.Author)
		name, err := tables.GetUsernameByID(id)
		if err != nil {
			http.Error(w, "Erreur récupération nom d'utilisateur", http.StatusInternalServerError)
		} else {
			p.Author = name
		}
		posts = append(posts, p)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
