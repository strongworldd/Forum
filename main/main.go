package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "forum/tables"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./BDD/accounts.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    http.HandleFunc("/api/register", registerHandler)
    http.HandleFunc("/api/login", loginHandler)
    http.Handle("/", http.FileServer(http.Dir("../")))

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