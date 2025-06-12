package tables

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

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