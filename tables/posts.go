package tables

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	//"strconv"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

// Fonction pour obtenir le chemin absolu vers le dossier BDD
func getDBPath() string {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	return filepath.Join(projectRoot, "BDD", "posts.db")
}

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
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de la base de données : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
	defer database.Close()

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, author TEXT)")
	if err != nil {
		log.Printf("Erreur lors de la préparation du statement : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Printf("Erreur lors de la création de la table posts : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
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

func CheckVotesTable() {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de la base de données : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
	defer database.Close()

	// Supprimer l'ancienne table si elle existe
	// database.Exec("DROP TABLE IF EXISTS votes") // Consider removing this line after initial setup

	// Créer la nouvelle table votes
	statement, err := database.Prepare(`
		CREATE TABLE IF NOT EXISTS votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			vote_type INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(post_id, username),
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Printf("Erreur lors de la préparation du statement votes : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
	defer statement.Close()

	_, err = statement.Exec()
	if err != nil {
		log.Printf("Erreur lors de la création de la table votes : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
}

// Fonction pour ajouter ou mettre à jour un vote
func AddOrUpdateVote(postID int, username string, voteType int) error {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer database.Close()

	// Vérifier si un vote existe déjà
	var existingVoteType int
	err = database.QueryRow("SELECT vote_type FROM votes WHERE post_id = ? AND username = ?", postID, username).Scan(&existingVoteType)

	if err == sql.ErrNoRows {
		// Pas de vote existant, on en crée un nouveau
		_, err = database.Exec("INSERT INTO votes (post_id, username, vote_type) VALUES (?, ?, ?)", postID, username, voteType)
	} else if err == nil {
		// Vote existant, on le met à jour
		_, err = database.Exec("UPDATE votes SET vote_type = ? WHERE post_id = ? AND username = ?", voteType, postID, username)
	}

	return err
}

// Fonction pour obtenir le vote d'un utilisateur sur un post
func GetUserVote(postID int, username string) (int, error) {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de l'ouverture de la base de données : %v", err)
	}
	defer database.Close()

	var voteType int
	err = database.QueryRow("SELECT vote_type FROM votes WHERE post_id = ? AND username = ?", postID, username).Scan(&voteType)
	if err == sql.ErrNoRows {
		return 0, nil // Pas de vote
	}
	return voteType, err
}

func InitDefaultPosts() {
	dbPath := getDBPath()
	database, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture de la base de données : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}
	defer database.Close()

	// Vérifier si des posts existent déjà
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	if err != nil {
		log.Printf("Erreur lors de la vérification des posts : %v", err)
		return // Ou gérer l'erreur de manière appropriée
	}

	// Si aucun post n'existe, ajouter les posts par défaut
	if count == 0 {
		defaultPosts := []struct {
			title   string
			content string
			author  string
		}{
			{
				title:   "Bienvenue sur notre forum !",
				content: "welcome.jpg",
				author:  "admin",
			},
			{
				title:   "Comment créer un post ?",
				content: "howto.jpg",
				author:  "admin",
			},
			{
				title:   "Les meilleures pratiques pour poster",
				content: "bestpractices.jpg",
				author:  "admin",
			},
			{
				title:   "Guide du débutant",
				content: "guide.jpg",
				author:  "admin",
			},
		}

		statement, err := database.Prepare("INSERT INTO posts (title, content, author) VALUES (?, ?, ?)")
		if err != nil {
			log.Printf("Erreur lors de la préparation du statement : %v", err)
			return // Ou gérer l'erreur de manière appropriée
		}
		defer statement.Close()

		for _, post := range defaultPosts {
			_, err = statement.Exec(post.title, post.content, post.author)
			if err != nil {
				log.Printf("Erreur lors de l'insertion du post %s : %v", post.title, err)
			}
		}
	}
}
