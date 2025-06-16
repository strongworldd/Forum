package models

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// Fonction pour obtenir le chemin absolu vers le dossier BDD
func getDBPath() string {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	return filepath.Join(projectRoot, "BDD", "posts.db")
}

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Created   time.Time `json:"created"`
	VoteCount int       `json:"voteCount"`
}

var (
	posts     = make(map[string]Post)
	postsLock sync.RWMutex
)

func AddPost(post Post) error {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer database.Close()

	// Vérifier si un post avec le même titre existe déjà
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM posts WHERE title = ?", post.Title).Scan(&count)
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification du post existant : %v", err)
	}
	if count > 0 {
		return ErrPostExists
	}

	// Insérer le nouveau post
	stmt, err := database.Prepare("INSERT INTO posts(title, content, author) VALUES(?, ?, ?)")
	if err != nil {
		return fmt.Errorf("erreur lors de la préparation de l'insertion du post : %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(post.Title, post.Content, post.Author)
	if err != nil {
		return fmt.Errorf("erreur lors de l'insertion du post : %v", err)
	}

	// Obtenir l'ID généré et l'assigner au post
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de l'ID du post : %v", err)
	}
	post.ID = int(id)

	log.Printf("Nouveau post ajouté à la base de données avec ID : %d", post.ID)

	return nil
}

func GetPost(title string) (Post, bool) {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de la base de données : %v", err)
		return Post{}, false
	}
	defer database.Close()

	var post Post
	row := database.QueryRow("SELECT id, title, content, author, created FROM posts WHERE title = ?", title)
	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.Created)
	if err == sql.ErrNoRows {
		return Post{}, false // Post non trouvé
	} else if err != nil {
		log.Printf("Erreur lors de la récupération du post : %v", err)
		return Post{}, false
	}

	// Récupérer le nombre de votes pour ce post
	voteCount, err := GetPostVoteCount(post.ID)
	if err != nil {
		log.Printf("Erreur lors de la récupération du nombre de votes pour le post %d : %v", post.ID, err)
		// Continuer même si les votes ne peuvent pas être récupérés
	}
	post.VoteCount = voteCount

	return post, true
}

func GetAllPosts() []Post {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de la base de données : %v", err)
		return nil
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title, content, author, created FROM posts ORDER BY created DESC")
	if err != nil {
		log.Printf("Erreur lors de la récupération de tous les posts : %v", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Author, &post.Created)
		if err != nil {
			log.Printf("Erreur lors du scan des données du post : %v", err)
			continue
		}

		// Récupérer le nombre de votes pour chaque post
		voteCount, err := GetPostVoteCount(post.ID)
		if err != nil {
			log.Printf("Erreur lors de la récupération du nombre de votes pour le post %d : %v", post.ID, err)
		}
		post.VoteCount = voteCount

		posts = append(posts, post)
	}

	return posts
}

func UpdatePostVoteCount(postID int, delta int) error {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer database.Close()

	// Récupérer le vote existant pour ce post (cela vient de la table votes, pas posts)
	// Note: La logique de vote est maintenant gérée dans tables/posts.go AddOrUpdateVote
	// Cette fonction devrait être mise à jour pour s'intégrer avec cela ou être supprimée si inutile.
	// Pour l'instant, on va juste s'assurer que cette fonction ne fait rien si elle est appelée.
	log.Printf("UpdatePostVoteCount appelé pour postID %d avec delta %d. Cette fonction devrait être refactorisée pour utiliser AddOrUpdateVote dans tables/posts.go.", postID, delta)
	return nil // Indiquer que le traitement est correct, même si rien n'est fait ici.
}

// Fonction pour obtenir le nombre total de votes d'un post
func GetPostVoteCount(postID int) (int, error) {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer database.Close()

	var totalVotes int
	err = database.QueryRow("SELECT COALESCE(SUM(vote_type), 0) FROM votes WHERE post_id = ?", postID).Scan(&totalVotes)
	return totalVotes, err
}

var (
	ErrPostExists   = &PostError{"post already exists"}
	ErrPostNotFound = &PostError{"post not found"}
)

type PostError struct {
	message string
}

func (e *PostError) Error() string {
	return e.message
}
