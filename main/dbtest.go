package main

import (
	"fmt"
	"forum/tables"
)

func Dbtest() {
	tables.CheckPostDB()
	tables.ResetPostsTable()
	tables.LoadAccounts()
	tables.LoadPosts()

	tables.CreatePost()
	tables.CreatePost()
	tables.CreatePost()

	tables.Deletepost(2)

	tables.LoadPosts()

	tables.CreatePost()

	tables.LoadPosts()

	fmt.Println("Database operations completed successfully.")
}
