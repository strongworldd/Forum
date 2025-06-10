package main

import (
	"fmt"
	"forum/tables"
)

func main() {
	tables.ResetAccountsTable()
	tables.LoadAccounts()
	tables.ResetPostsTable()
	fmt.Println("Database operations completed successfully.")
}
