package main

import (
	"database/sql"
	"fmt"
	"forum"
	"os"
	"github.com/mattn/go-sqlite3"
)

// Removed self-referential import

func main() {
	fileName := "./bdd.db"
	
	_, error := os.Stat(fileName)

    // check if error is "file not exists"
    if os.IsNotExist(error) {
    	fmt.Printf("%v file does not exist\n", fileName)
		os.Create(fileName)
		CREATE TABLE IF NOT EXISTS product(product_id int primary key auto_increment, product_name text, product_price int, created_at datetime default CURRENT_TIMESTAMP, updated_at datetime default CURRENT_TIMESTAMP)

	} else {
		db, err := sql.Open("sqlite3",fileName)
		if err != nil {
				fmt.Println("Error opening database:", err)
			return
		}
		result, err := db.Query("SELECT name, age FROM User WHERE age > 30")
		var name string
		var age int
		for result.Next() {
		result.Scan(&name, &age)
		/* Faire quelque chose avec cette ligne */
		}
		result.Close()
	}
}