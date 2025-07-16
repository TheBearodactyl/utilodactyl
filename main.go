package main

import (
	"utilodactyl/actions"
)

func main() {
	err := actions.App()
	if err != nil {
		panic(err)
	}
}
