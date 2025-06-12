package main

import (
	"fmt"
	"forum/tables"
)

func main() {
	tables.CheckPostDB()
	tables.CheckAccountDB()

	tables.ResetPostsTable()
	tables.ResetAccountsTable()

	tables.LoadAccounts()

	tables.CreateAccount("Alice", "Alice@gmail.com", "password123")
	tables.CreateAccount("Bob", "Bob@gmail.com", "password123")
	tables.CreateAccount("John Doe", "JohnDoe@gmail.com", "password123")
	tables.CreateAccount("Kerchak", "Kerchak@gmail.com", "password123")
	tables.CheckAccountName("Josh")

	tables.LoadAccounts()

	fmt.Println(tables.CheckConnexion("Caca", "password123"))
	fmt.Println(tables.CheckConnexion("Alice", "password123"))
	fmt.Println(tables.CheckConnexion("Kerchak@gmail.com", "password123"))

	fmt.Println("Database operations completed successfully.")
}
