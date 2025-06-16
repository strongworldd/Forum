package tables

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"encoding/json"
	"strconv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(name, hashedPassword, email string) error {
	stmt, err := repo.DB.Prepare(`INSERT INTO people (name, password, email, postliked) VALUES (?, ?, ?, ?)`)
    if err != nil {
        return fmt.Errorf("failed to prepare insert user statement: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(name, hashedPassword, email, "")
    if err != nil {
        return fmt.Errorf("failed to execute insert user statement: %w", err)
    }
    return nil
}

func (repo *UserRepository) GetPasswordByUsername(name string) (string, error) {
	var hashedPassword string

	row := repo.DB.QueryRow(`SELECT password FROM people WHERE name = ?`, name)
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil 
		}
		return "", fmt.Errorf("failed to query password: %w", err)
	}
	return hashedPassword, nil
}

func (repo *UserRepository) GetUserID(name string) (int, error) {
	var userID int

	row := repo.DB.QueryRow(`SELECT id FROM people WHERE name = ?`, name)
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil 
		}
		return 0, fmt.Errorf("failed to query user id: %w", err)
	}
	return userID, nil
}

func GetUserByID(id int) (*User, error) {
	database, _ := sql.Open("sqlite3", "../BDD/accounts.db")
    defer database.Close()

    row := database.QueryRow("SELECT id, name, email FROM people WHERE id = ?", id)
    var user User
    err := row.Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
		fmt.Println("Error fetching user:", err)
        return nil, err
    }
    return &user, nil
}

func GetUsernameByID(id int) (string, error) {
	user, err := GetUserByID(id)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return "", err
	}
	return user.Username, nil
}

func (repo *UserRepository) ListAllUsernames() ([]string, error) {
	rows, err := repo.DB.Query(`SELECT name FROM people`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all usernames: %w", err)
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan username: %w", err)
		}
		usernames = append(usernames, name)
	}
	return usernames, nil
}

func CheckUsernameExists(name string) (bool, error) {
	database, _ := sql.Open("sqlite3", "../BDD/accounts.db")
	defer database.Close()

	row := database.QueryRow("SELECT COUNT(*) FROM people WHERE name = ?", name)
	var count int
	err := row.Scan(&count)
	fmt.Println("Checking if username exists:", name, "Count:", count)
	if err != nil {
		return false, fmt.Errorf("failed to check if username exists: %w", err)
	}
	return count > 0, nil
}

func CheckEmailExists(email string) (bool, error) {
	database, _ := sql.Open("sqlite3", "../BDD/accounts.db")
	defer database.Close()

	row := database.QueryRow("SELECT COUNT(*) FROM people WHERE email = ?", email)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}
	return count > 0, nil
}

func ComparePasswordsWithId(id int, password string) (bool, error) {
	database, _ := sql.Open("sqlite3", "../BDD/accounts.db")
	defer database.Close()

	row := database.QueryRow("SELECT password FROM people WHERE id = ?", id)
	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		return false, fmt.Errorf("failed to query password by id: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, fmt.Errorf("password comparison failed: %w", err)
	}
	return true, nil
}

func SaveSettingsToUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	var data struct {
		Username 	string `json:"username"`
		Email 	 	string `json:"email"`
		Newpassword string `json:"newpassword"`
		Password 	string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&data)

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

	correct, err := ComparePasswordsWithId(userID, data.Password)
	if err != nil || !correct {
		http.Error(w, "Mot de passe incorrect", http.StatusUnauthorized)
		return
	}

	database, _ := sql.Open("sqlite3", "../BDD/accounts.db")
    defer database.Close()

	if data.Username != "" {
		username, err := CheckUsernameExists(data.Username)
		row := database.QueryRow("SELECT name FROM people WHERE id = ?", userID)
		var actualUsername string
		_ = row.Scan(&actualUsername)
		if actualUsername == data.Username {
			err = nil
			username = false
		}

		if (err != nil || username) {
			http.Error(w, "failed to check username existence: %w", http.StatusUnauthorized)
			return
		}
		_, err = database.Exec("UPDATE people SET name = ? WHERE id = ?", data.Username, userID)
		if err != nil {
			http.Error(w, "failed to update username: %w", http.StatusInternalServerError)
			return
		}
	}

	if data.Email != "" {
		mail, err := CheckEmailExists(data.Email)
		row := database.QueryRow("SELECT email FROM people WHERE id = ?", userID)
		var actualEmail string
		_ = row.Scan(&actualEmail)
		if actualEmail == data.Email {
			err = nil
			mail = false
		}

		if err != nil || mail {
			http.Error(w, "failed to check email existence: %w", http.StatusUnauthorized)
			return
		}
		_, err = database.Exec("UPDATE people SET email = ? WHERE id = ?", data.Email, userID)
		if err != nil {
			http.Error(w, "failed to update mail: %w", http.StatusInternalServerError)
			return
		}
	}

	if data.Newpassword != "" {
		newpassword, _ := HashPassword(data.Newpassword)
		_, err = database.Exec("UPDATE people SET password = ? WHERE id = ?", newpassword, userID)
		if err != nil {
			http.Error(w, "failed to update password: %w", http.StatusInternalServerError)
			return
		}
	}
}

func (repo *UserRepository) GetUsernameByUUID(userUUID uuid.UUID) (string, error) {
	var name string

	row := repo.DB.QueryRow(`SELECT name FROM people WHERE uuid = ?`, userUUID.String())
	err := row.Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil 
		}
		return "", fmt.Errorf("failed to query username by uuid: %w", err)
	}
	return name, nil
}

func (repo *UserRepository) UpdateUUIDForUser(userUUID uuid.UUID, name string) error {
	stmt, err := repo.DB.Prepare(`UPDATE people SET uuid = ? WHERE name = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update uuid statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userUUID.String(), name)
	if err != nil {
		return fmt.Errorf("failed to execute update uuid statement: %w", err)
	}
	return nil
}
