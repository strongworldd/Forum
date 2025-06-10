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

func LoadAccounts() {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, name TEXT, password TEXT, postliked TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (name, password, postliked) VALUES (?, ?, ?)")
	statement.Exec("John Doe", "secret123", arraytostr([]int{1, 2, 3}))
	rows, _ := database.Query("SELECT id, name, password, postliked FROM people")
	fmt.Println(rows)
	var id int
	var name string
	var password string
	var postliked string
	for rows.Next() {
		rows.Scan(&id, &name, &password, &postliked)
		fmt.Println(strconv.Itoa(id) + ": " + name + " pass: " + password + " Likes: " + postliked)
	}
	fmt.Println("End")
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