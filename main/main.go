package main

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"Forum/handlers"
	"Forum/tables"
)

func main() {
	// Initialiser les tables
	tables.CheckPostDB()
	tables.CheckVotesTable()
	tables.InitDefaultPosts()

	// Routes pour l'authentification
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)

	// Routes pour les posts
	http.HandleFunc("/api/posts", handlers.GetPostsHandler)
	http.HandleFunc("/api/posts/create", handlers.CreatePostHandler)

	// Routes pour les votes
	http.HandleFunc("/api/posts/vote", handlers.VoteHandler)
	http.HandleFunc("/api/posts/votes", handlers.GetUserVotesHandler)

	// Obtenir le chemin absolu vers le dossier racine du projet
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))

	// Servir les fichiers statiques
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(projectRoot, "css")))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(filepath.Join(projectRoot, "img")))))
	http.Handle("/Login/", http.StripPrefix("/Login/", http.FileServer(http.Dir(filepath.Join(projectRoot, "Login")))))
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir(filepath.Join(projectRoot, "html")))))

	// Rediriger la racine vers home.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/html/home.html", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// Démarrer le serveur
	log.Println("Serveur démarré sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
