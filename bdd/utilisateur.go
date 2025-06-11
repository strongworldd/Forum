package BDD

import (
	"database/sql"
	"fmt"

	uuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func NewUser(Pseudo string, HashPass string, db *sql.DB) int {
	statement, err := db.Prepare("INSERT INTO User (pseudo, password) VALUES(?,?)")
	if err != nil {
		fmt.Println(err)
		fmt.Println("error Prepare new user")
		return (500)
	}
	statement.Exec(Pseudo, HashPass)
	db.Close()
	return (0)
}

func CheckPassword(username string, db *sql.DB) (int, string) {
	var HashPass string
	var password string

	tsql, err := db.Query("SELECT password FROM User WHERE pseudo = (?)", username)
	if err != nil {
		fmt.Println(err)
		return 500, HashPass
	}

	for tsql.Next() {
		tsql.Scan(&password)
	}
	HashPass = password
	return 0, HashPass
}

func GetId_User(username string, db *sql.DB) (int, int) {
	var Id_user int

	tsql, err := db.Query("SELECT Id_user FROM User WHERE pseudo = (?)", username)
	if err != nil {
		fmt.Println(err)
		return 500, Id_user
	}

	for tsql.Next() {
		tsql.Scan(&Id_user)
	}
	if Id_user > 0 {
		return 0, Id_user
	}
	return 500, Id_user
}

func GetAllUsername(db *sql.DB) (int, []string) {
	var allUsername []string
	var username string

	tsql, err := db.Query("SELECT pseudo FROM User")
	if err != nil {
		fmt.Println(err)
		return 500, allUsername
	}

	for tsql.Next() {
		tsql.Scan(&username)
		allUsername = append(allUsername, username)
	}
	return 0, allUsername
}

func GetUserByUUID(uuid string, db *sql.DB) (int, string) {
	var username string

	tsql, err := db.Query("SELECT pseudo FROM User WHERE uuid = (?)", uuid) // check for UUID name in database
	if err != nil {
		fmt.Println(err)
		return 400, username
	}
	for tsql.Next() {
		tsql.Scan(&username)
	}
	return 0, username
}

func PutUUID(UUID uuid.UUID, pseudo string, db *sql.DB) int {
	statement, err := db.Prepare("UPDATE User SET uuid = ? WHERE (pseudo =?)")
	if err != nil {
		fmt.Println(err)
		fmt.Println("error Prepare new user")
		return (500)
	}
	statement.Exec(UUID, pseudo)
	db.Close()
	return (0)
}