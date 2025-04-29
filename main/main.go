package main

import "forum"

// Removed self-referential import

func main() {
	forum.Delete("test.txt")
}