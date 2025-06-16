package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"Forum/models"
	"github.com/google/uuid"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Parse le formulaire multipart avec une limite de 10 MB pour les fichiers
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'analyse du formulaire : %v", err), http.StatusBadRequest)
		return
	}

	// Récupérer le titre et l'auteur du formulaire
	title := r.FormValue("title")
	author := r.FormValue("author")
	// L'auteur peut être récupéré de l'en-tête X-Username si la session est gérée
	// Pour l'instant, utilisons celui du formulaire pour la simplicité.
	if author == "" {
		http.Error(w, "Le champ auteur est requis", http.StatusBadRequest)
		return
	}
	if title == "" {
		http.Error(w, "Le champ titre est requis", http.StatusBadRequest)
		return
	}

	// Récupérer le fichier image
	file, fileHeader, err := r.FormFile("content") // "content" est le nom du champ input type="file"
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la récupération du fichier : %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Générer un nom de fichier unique
	fileExt := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + fileExt

	// Obtenir le chemin absolu vers le dossier img
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	imgPath := filepath.Join(projectRoot, "img", newFileName)

	// Créer le dossier img si il n'existe pas
	if _, err := os.Stat(filepath.Join(projectRoot, "img")); os.IsNotExist(err) {
		_ = os.Mkdir(filepath.Join(projectRoot, "img"), 0755)
	}

	// Créer un nouveau fichier sur le serveur pour y copier l'image
	dst, err := os.Create(imgPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la création du fichier image sur le serveur : %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copier le contenu du fichier uploadé vers le nouveau fichier
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de la copie du fichier image : %v", err), http.StatusInternalServerError)
		return
	}

	// Créer le post dans la base de données
	post := models.Post{
		Title:   title,
		Content: newFileName, // Le nom du fichier image est stocké comme contenu
		Author:  author,
	}

	if err := models.AddPost(post); err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'ajout du post : %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	posts := models.GetAllPosts()
	json.NewEncoder(w).Encode(posts)
}

// ServeStaticFiles n'est plus nécessaire ici car géré dans main.go
func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r) // S'assurer que cette fonction ne fait rien si elle est accidentellement appelée.
}
