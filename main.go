package main

import (
	"siteutil/actions"
)

func main() {
	err := actions.App()
	if err != nil {
		panic(err)
	}
}
