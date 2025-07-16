package main

import (
	"utilodactyl/actions"
	"utilodactyl/models"

	"github.com/alexflint/go-arg"
)

func main() {
	arg.MustParse(&models.Cli)
	err := actions.App()
	if err != nil {
		panic(err)
	}
}
