package tables

import (
	"database/sql"
	"fmt"
)

func CheckConnexion(identifier string, password string) int {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()
	rows, _ := database.Query("SELECT id FROM people WHERE (name = ? OR mail = ?) AND password = ?", identifier, identifier, password)
	defer rows.Close()

	var id int
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return 0
		}
	}

	return id
}