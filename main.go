package main

import (
	"fmt"
	"utilodactyl/actions"
	"utilodactyl/models"

	"github.com/alexflint/go-arg"
)

func main() {
	fmt.Print("\033[H\033[2J")
	arg.MustParse(&models.Cli)
	err := actions.App()
	if err != nil {
		panic(err)
	}
}
