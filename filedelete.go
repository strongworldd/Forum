package main

import (
	"fmt"
	"os"
)

func Delete(file string) {
	// Delete the file
	err := os.Remove(file)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	fmt.Println("File deleted successfully")
}