package main

import (
	"fmt"
	"os"
)

func Create(name string) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	l, err := f.WriteString("Hello World")
	if err != nil {
		f.Close()
		fmt.Println("Error writing on file:", err)
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
		return
	}

	fmt.Println("l:", l)
}