package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"Forum/models"
	"Forum/tables"

	_ "github.com/mattn/go-sqlite3"
)

// var db *sql.DB // Removed global db variable

// func init() { // Removed init function
// 	// Obtenir le chemin absolu vers le fichier de base de données
// 	_, currentFile, _, _ := runtime.Caller(0)
// 	projectRoot := filepath.Dir(filepath.Dir(currentFile))
// 	dbPath := filepath.Join(projectRoot, "BDD", "posts.db")

// 	var err error
// 	db, err = sql.Open("sqlite3", dbPath)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username := r.Header.Get("X-Username")
	if username == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	var data struct {
		Title    string `json:"title"`
		Author   string `json:"author"`
		VoteType int    `json:"voteType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erreur de décodage", http.StatusBadRequest)
		return
	}

	// Récupérer l'ID du post
	// Use models.GetPost to get the post and its ID
	post, exists := models.GetPost(data.Title)
	if !exists {
		http.Error(w, "Post non trouvé", http.StatusNotFound)
		return
	}

	// Ajouter ou mettre à jour le vote
	if err := tables.AddOrUpdateVote(post.ID, username, data.VoteType); err != nil {
		http.Error(w, "Erreur lors du vote", http.StatusInternalServerError)
		return
	}

	// Récupérer le nouveau nombre total de votes en utilisant models.GetPostVoteCount
	totalVotes, err := models.GetPostVoteCount(post.ID)
	if err != nil {
		http.Error(w, "Erreur lors du calcul des votes", http.StatusInternalServerError)
		return
	}

	// Récupérer le vote de l'utilisateur en utilisant tables.GetUserVote
	userVote, err := tables.GetUserVote(post.ID, username)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération du vote", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"voteCount": totalVotes,
		"voteType":  userVote,
	})
}

func GetUserVotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username := r.Header.Get("X-Username")
	if username == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	// Récupérer tous les posts avec leurs votes
	// Le code ici est un peu plus complexe car il ne dépendait pas directement du `db` global.
	// Il faut s'assurer que les fonctions `tables.GetUserVote` sont utilisées correctement ici.
	// Pour l'instant, laissons-le tel quel, mais il faudra revoir cette partie si des problèmes persistent.

	// Example of refactoring GetUserVotesHandler to use existing models/tables functions
	posts := models.GetAllPosts()                        // Get all posts from the database
	votesData := make(map[string]map[string]interface{}) // Map[title] -> {count, userVote}

	for _, post := range posts {
		// Get total votes for the post
		totalVotes, err := models.GetPostVoteCount(post.ID)
		if err != nil {
			log.Printf("Erreur lors de la récupération du nombre total de votes pour le post %d: %v", post.ID, err)
			totalVotes = 0 // Default to 0 on error
		}

		// Get user's vote for the post
		userVote, err := tables.GetUserVote(post.ID, username)
		if err != nil {
			log.Printf("Erreur lors de la récupération du vote de l'utilisateur pour le post %d: %v", post.ID, err)
			userVote = 0 // Default to 0 on error
		}

		votesData[post.Title] = map[string]interface{}{
			"count":    totalVotes,
			"userVote": userVote,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votesData)
}
