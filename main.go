package main

import (
	"bookify/actions"
	"fmt"
	"log"
)

func main() {
	if err := actions.App(); err != nil {
		log.Fatalf("Application error: %v", err)
	}

	fmt.Println("Exiting Bookify. Goodbye!")
}
