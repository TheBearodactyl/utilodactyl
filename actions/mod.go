package actions

import (
	"bookify/actions/add"
	"bookify/actions/edit"
	"bookify/actions/update"
	"bookify/actions/view"
	"bookify/actions/vue"
	"fmt"
	"github.com/charmbracelet/huh"
)

type AppAction string

const (
	AddBook    AppAction = "Add a new book"
	ViewBooks  AppAction = "View existing books"
	EditBook   AppAction = "Edit a book"
	GenVue     AppAction = "Generate Vue component data"
	UpdateData AppAction = "Update the `books.json` release"
	ExitApp    AppAction = "Exit"
)

func App() error {
	for {
		var action AppAction
		err := huh.NewSelect[AppAction]().
			Title("What would you like to do?").
			Options(
				huh.NewOption(string(AddBook), AddBook),
				huh.NewOption(string(ViewBooks), ViewBooks),
				huh.NewOption(string(EditBook), EditBook),
				huh.NewOption(string(GenVue), GenVue),
				huh.NewOption(string(UpdateData), UpdateData),
				huh.NewOption(string(ExitApp), ExitApp),
			).
			Value(&action).
			Run()

		if err != nil {
			return fmt.Errorf("prompt error: %w", err)
		}

		switch action {
		case AddBook:
			if err := add.AddBook(); err != nil {
				fmt.Printf("Error adding book: %v\n", err)
			}
		case ViewBooks:
			if err := view.ViewBooks(); err != nil {
				fmt.Printf("Error viewing books: %v\n", err)
			}
		case EditBook:
			if err := edit.EditBook(); err != nil {
				fmt.Printf("Error editing book: %v\n", err)
			}
		case GenVue:
			if err := vue.GenerateVueData(); err != nil {
				fmt.Printf("Error generating Vue data: %v\n", err)
			}
		case UpdateData:
			if err := update.UpdateBooks(); err != nil {
				fmt.Printf("Error updating books.json: %v\n", err)
			}
		case ExitApp:
			return nil
		}
	}
}
