package tables

import (
	"database/sql"
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

func ResetAccountsTable() {
    database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
    defer database.Close()
    _, _ = database.Exec("DELETE FROM people")
}

func DeleteAccount(id int) {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
    defer database.Close()
    _, _ = database.Exec("DELETE FROM people WHERE id = ?", id)
}

func CheckAccountDB() {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT, mail TEXT, password TEXT, postliked TEXT)")
	statement.Exec()
}

func CheckAccountName(username string) bool {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	rows, _ := database.Query("SELECT id, name, password, postliked FROM people WHERE name = ?", username)
	defer rows.Close()

	return rows.Next()
}

func CreateAccount(username string, mail string, password string) bool{
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	if !CheckAccountName(username) {
		statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT, mail TEXT, password TEXT, postliked TEXT)")
		statement.Exec()

		statement, _ = database.Prepare("INSERT INTO people (name, mail, password, postliked) VALUES (?, ?, ?, ?)")
		statement.Exec(username, mail, password, "")
		fmt.Println("Account created successfully for:", username)
		return true
	}
	return false
}

func LoadAccounts() {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	rows, _ := database.Query("SELECT id, name, mail, password, postliked FROM people")
	defer rows.Close()

	var id int
	var name string
	var mail string
	var password string
	var postliked string
	for rows.Next() {
		rows.Scan(&id, &name, &mail, &password, &postliked)
		fmt.Println(strconv.Itoa(id) + ": " + name + " email: " + mail + " pass: " + password + " Likes: " + postliked)
	}
}

func strtoarray(s string) []int {
	var arr []int
	var temp = ""
	for _, v := range s {
		if s[v] == ',' {
			if temp != "" {
				num, err := strconv.Atoi(temp)
				if err == nil {
					arr = append(arr, num)
				}
				temp = ""
			}
		} else {
			temp += string(s[v])
		}
	}
	return arr
}

func arraytostr(arr []int) string {
	var str string
	for i, v := range arr {
		str += strconv.Itoa(v)
		if i < len(arr)-1 {
			str += ","
		}
	}
	return str
}