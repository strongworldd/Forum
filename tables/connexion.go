package tables

import (
	"database/sql"
	"fmt"
)

func CheckConnexion(identifier string, password string) int {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()
	
	rows, _ := database.Query("SELECT id, password FROM people WHERE (name = ? OR mail = ?)", identifier, identifier)
	defer rows.Close()

	var id int
	var storedPassword string
	if rows.Next() {
		err := rows.Scan(&id, &storedPassword)
		if !checkPasswordHash(password, storedPassword) {
			fmt.Println("Invalid password for user ID:", id)
			return 0
		}
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return 0
		}
	}

	return id
}