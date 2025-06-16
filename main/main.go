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
	http.HandleFunc("/api/posts/delete", deletePostHandler)
	http.HandleFunc("/api/posts/edit", editPostHandler)
	http.HandleFunc("/api/posts/vote", voteHandler)
	http.HandleFunc("/api/posts/votes", getVotesHandler)

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

	// On récupère le fichier image
	file, handler, err := r.FormFile("content")
	if err != nil {
		http.Error(w, "Erreur fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// On sauve le fichier dans le dossier img/
	imgPath := "../img/" + handler.Filename
	dst, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, "Erreur sauvegarde image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

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

	_, err = stmt.Exec(title, handler.Filename, author)
	if err != nil {
		http.Error(w, "Erreur insertion", 500)
		return
	}

	http.Redirect(w, r, "/html/home.html", http.StatusSeeOther)
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

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Title  string `json:"title"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erreur de décodage", http.StatusBadRequest)
		return
	}

	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Supprimer le post
	stmt, err := dbPosts.Prepare("DELETE FROM posts WHERE title = ? AND author = ?")
	if err != nil {
		http.Error(w, "Erreur préparation", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(data.Title, data.Author)
	if err != nil {
		http.Error(w, "Erreur suppression", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Erreur vérification", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post supprimé avec succès"})
}

func editPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		OldTitle   string `json:"oldTitle"`
		OldAuthor  string `json:"oldAuthor"`
		NewTitle   string `json:"newTitle"`
		NewContent string `json:"newContent"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erreur de décodage", http.StatusBadRequest)
		return
	}

	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Mettre à jour le post
	stmt, err := dbPosts.Prepare("UPDATE posts SET title = ?, content = ? WHERE title = ? AND author = ?")
	if err != nil {
		http.Error(w, "Erreur préparation", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(data.NewTitle, data.NewContent, data.OldTitle, data.OldAuthor)
	if err != nil {
		http.Error(w, "Erreur mise à jour", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Erreur vérification", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post modifié avec succès"})
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Vérifier si l'utilisateur est connecté
	username := r.Header.Get("X-Username")
	if username == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	var data struct {
		Title    string `json:"title"`
		Author   string `json:"author"`
		VoteType int    `json:"voteType"` // 1 pour upvote, -1 pour downvote
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erreur de décodage", http.StatusBadRequest)
		return
	}

	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Récupérer l'ID du post
	var postID int
	err = dbPosts.QueryRow("SELECT id FROM posts WHERE title = ? AND author = ?", data.Title, data.Author).Scan(&postID)
	if err != nil {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	// Mettre à jour ou insérer le vote
	stmt, err := dbPosts.Prepare(`
		INSERT INTO votes (post_id, username, vote_type) 
		VALUES (?, ?, ?) 
		ON CONFLICT(post_id, username) 
		DO UPDATE SET vote_type = excluded.vote_type
	`)
	if err != nil {
		http.Error(w, "Erreur préparation", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(postID, username, data.VoteType)
	if err != nil {
		http.Error(w, "Erreur vote", http.StatusInternalServerError)
		return
	}

	// Récupérer le nombre total de votes
	var totalVotes int
	err = dbPosts.QueryRow(`
		SELECT COALESCE(SUM(vote_type), 0)
		FROM votes
		WHERE post_id = ?
	`, postID).Scan(&totalVotes)
	if err != nil {
		http.Error(w, "Erreur calcul votes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Vote enregistré",
		"voteType":  data.VoteType,
		"voteCount": totalVotes,
	})
}

func getVotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username := r.Header.Get("X-Username")
	if username == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	dbPosts, err := sql.Open("sqlite3", "../BDD/posts.db")
	if err != nil {
		http.Error(w, "Erreur BDD", http.StatusInternalServerError)
		return
	}
	defer dbPosts.Close()

	// Récupérer tous les posts avec leurs votes
	rows, err := dbPosts.Query(`
		SELECT p.title, 
		       COALESCE(SUM(v.vote_type), 0) as total_votes,
		       (SELECT vote_type FROM votes WHERE post_id = p.id AND username = ?) as user_vote
		FROM posts p
		LEFT JOIN votes v ON p.id = v.post_id
		GROUP BY p.id, p.title
	`, username)
	if err != nil {
		http.Error(w, "Erreur lecture votes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type VoteData struct {
		Count    int `json:"count"`
		UserVote int `json:"userVote"`
	}
	votes := make(map[string]VoteData)

	for rows.Next() {
		var title string
		var totalVotes int
		var userVote sql.NullInt64
		if err := rows.Scan(&title, &totalVotes, &userVote); err != nil {
			continue
		}
		votes[title] = VoteData{
			Count:    totalVotes,
			UserVote: int(userVote.Int64),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votes)
}
