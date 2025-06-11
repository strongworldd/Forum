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

	tables.CreateAccount("Alice", "password123")
	tables.CheckAccountName("Josh")

	tables.LoadAccounts()

	fmt.Println("Database operations completed successfully.")
}
