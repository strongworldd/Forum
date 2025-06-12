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

func (repo *UserRepository) CreateUser(pseudo, hashedPassword string) error {
	stmt, err := repo.DB.Prepare(`INSERT INTO User (pseudo, password) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert user statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(pseudo, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to execute insert user statement: %w", err)
	}
	return nil
}

func (repo *UserRepository) GetPasswordByUsername(username string) (string, error) {
	var hashedPassword string

	row := repo.DB.QueryRow(`SELECT password FROM User WHERE pseudo = ?`, username)
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // user not found
		}
		return "", fmt.Errorf("failed to query password: %w", err)
	}
	return hashedPassword, nil
}

func (repo *UserRepository) GetUserID(username string) (int, error) {
	var userID int

	row := repo.DB.QueryRow(`SELECT Id_user FROM User WHERE pseudo = ?`, username)
	err := row.Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil // user not found
		}
		return 0, fmt.Errorf("failed to query user id: %w", err)
	}
	return userID, nil
}

func (repo *UserRepository) ListAllUsernames() ([]string, error) {
	rows, err := repo.DB.Query(`SELECT pseudo FROM User`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all usernames: %w", err)
	}
	defer rows.Close()

	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to scan username: %w", err)
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}

func (repo *UserRepository) GetUsernameByUUID(userUUID uuid.UUID) (string, error) {
	var username string

	row := repo.DB.QueryRow(`SELECT pseudo FROM User WHERE uuid = ?`, userUUID.String())
	err := row.Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil // not found
		}
		return "", fmt.Errorf("failed to query username by uuid: %w", err)
	}
	return username, nil
}

func (repo *UserRepository) UpdateUUIDForUser(userUUID uuid.UUID, pseudo string) error {
	stmt, err := repo.DB.Prepare(`UPDATE User SET uuid = ? WHERE pseudo = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update uuid statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userUUID.String(), pseudo)
	if err != nil {
		return fmt.Errorf("failed to execute update uuid statement: %w", err)
	}
	return nil
}