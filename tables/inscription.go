package tables

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type AccountRepository struct {
	DB *sql.DB
}

func OpenDatabase(filepath string) (*AccountRepository, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &AccountRepository{DB: db}, nil
}

func (repo *AccountRepository) InitAccountTable() error {
    _, err := repo.DB.Exec(`
        CREATE TABLE IF NOT EXISTS people (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT UNIQUE,
            password TEXT,
            postliked TEXT,
            uuid TEXT
        )
    `)
    return err
}

func (repo *AccountRepository) RegisterUser(username, password string) error {
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}

	_, err = repo.DB.Exec(`
		INSERT INTO people (name, password, postliked) VALUES (?, ?, '')
	`, username, hashed)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (repo *AccountRepository) AuthenticateUser(username, password string) error {
	var storedHash string
	err := repo.DB.QueryRow(`SELECT password FROM people WHERE name = ?`, username).Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func RenderSignupPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"Forum/html/home.html",
		"Forum/html/sidebar.html",
	)
	if err != nil {
		log.Printf("error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}