package tables

import (
	"database/sql"
	"fmt"
	"strconv"
	_ "github.com/mattn/go-sqlite3" 
	"golang.org/x/crypto/bcrypt"
	"log"
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

	// On uniformise bien la structure ici :
	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS people (
		id INTEGER PRIMARY KEY, 
		name TEXT, 
		password TEXT, 
		email TEXT, 
		postliked TEXT
	)`)
	statement.Exec()
}

func CreateAccount(username string, password string) {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println("Erreur lors du hash:", err)
		return
	}

	statement, _ := database.Prepare("INSERT INTO people (name, password, email, postliked) VALUES (?, ?, ?, ?)")
	statement.Exec(username, hashedPassword, "test@email.com", "")
}

func LoadAccounts() {
	database, _ := sql.Open("sqlite3", "./BDD/accounts.db")
	defer database.Close()

	rows, _ := database.Query("SELECT id, name, password, email, postliked FROM people")
	defer rows.Close()
	
	var id int
	var name string
	var password string
	var email string
	var postliked string
	for rows.Next() {
		rows.Scan(&id, &name, &password, &email, &postliked)
		fmt.Println(strconv.Itoa(id) + ": " + name + " pass: " + password + " email: " + email + " Likes: " + postliked)
	}
}

func Strtoarray(s string) []int {
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

func Arraytostr(arr []int) string {
	var str string
	for i, v := range arr {
		str += strconv.Itoa(v)
		if i < len(arr)-1 {
			str += ","
		}
	}
	return str
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func ComparePasswords(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("Erreur de comparaison de mot de passe:", err)
		return false
	}
	return true
}
