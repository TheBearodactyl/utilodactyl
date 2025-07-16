// Package actions
package actions

import (
	"fmt"
	"github.com/charmbracelet/huh"
	bookadd "siteutil/actions/books/add"
	bookedit "siteutil/actions/books/edit"
	bookupdate "siteutil/actions/books/update"
	bookview "siteutil/actions/books/view"
	gameadd "siteutil/actions/games/add"
	gameedit "siteutil/actions/games/edit"
	gameupdate "siteutil/actions/games/update"
	gameview "siteutil/actions/games/view"
	projectadd "siteutil/actions/projects/add"
	projectedit "siteutil/actions/projects/edit"
	projectupdate "siteutil/actions/projects/update"
	projectview "siteutil/actions/projects/view"
)

type AppAction string

const (
	Books          AppAction = "Operate on `books.json`"
	Projects       AppAction = "Operate on `projects.json`"
	Games          AppAction = "Operate on `games.json`"
	AddBook        AppAction = "Add a new book"
	AddProject     AppAction = "Add a new project"
	AddGame        AppAction = "Add a new game"
	ViewBooks      AppAction = "View existing books"
	ViewProject    AppAction = "View existing projects"
	ViewGames      AppAction = "View existing games"
	EditBook       AppAction = "Edit a book"
	EditProject    AppAction = "Edit a project"
	EditGame       AppAction = "Edit a game"
	UpdateBooks    AppAction = "Update the `books.json` release"
	UpdateProjects AppAction = "Update the `projects.json` release"
	UpdateGames    AppAction = "Update the `games.json` release"
	ExitApp        AppAction = "Exit"
)

func App() error {
	for {
		var action AppAction
		err := huh.NewSelect[AppAction]().
			Title("What would you like to do?").
			Options(
				huh.NewOption(string(Books), Books),
				huh.NewOption(string(Projects), Projects),
				huh.NewOption(string(Games), Games),
			).
			Value(&action).Run()

		if err != nil {
			return fmt.Errorf("%w", err)
		}

		switch action {
		case Books:
			var bookAction AppAction
			err := huh.NewSelect[AppAction]().
				Title("What would you like to do?").
				Options(
					huh.NewOption(string(AddBook), AddBook),
					huh.NewOption(string(ViewBooks), ViewBooks),
					huh.NewOption(string(EditBook), EditBook),
					huh.NewOption(string(UpdateBooks), UpdateBooks),
					huh.NewOption(string(ExitApp), ExitApp),
				).
				Value(&bookAction).
				Run()

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			switch bookAction {
			case AddBook:
				if err := bookadd.AddBook(); err != nil {
					fmt.Printf("Error adding book: %v\n", err)
				}
			case ViewBooks:
				if err := bookview.ViewBooks(); err != nil {
					fmt.Printf("Error viewing books: %v\n", err)
				}
			case EditBook:
				if err := bookedit.EditBook(); err != nil {
					fmt.Printf("Error editing book: %v\n", err)
				}
			case UpdateBooks:
				if err := bookupdate.UpdateBooks(); err != nil {
					fmt.Printf("Error updating books: %v\n", err)
				}
			}
		case Projects:
			var projectAction AppAction
			err := huh.NewSelect[AppAction]().
				Title("What would you like to do?").
				Options(
					huh.NewOption(string(AddProject), AddProject),
					huh.NewOption(string(EditProject), EditProject),
					huh.NewOption(string(ViewProject), ViewProject),
					huh.NewOption(string(UpdateProjects), UpdateProjects),
				).
				Value(&projectAction).
				Run()

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			switch projectAction {
			case AddProject:
				if err := projectadd.AddProject(); err != nil {
					fmt.Printf("Error adding project: %v\n", err)
				}
			case EditProject:
				if err := projectedit.EditProject(); err != nil {
					fmt.Printf("Error editing project: %v\n", err)
				}
			case ViewProject:
				if err := projectview.ViewProjects(); err != nil {
					fmt.Printf("Error viewing projects: %v\n", err)
				}
			case UpdateProjects:
				if err := projectupdate.UpdateProjects(); err != nil {
					fmt.Printf("Error updating projects: %v\n", err)
				}
			}
		case Games:
			var gamesAction AppAction
			err := huh.NewSelect[AppAction]().
				Title("What would you like to do?").
				Options(
					huh.NewOption(string(AddGame), AddGame),
					huh.NewOption(string(EditGame), EditGame),
					huh.NewOption(string(UpdateGames), UpdateGames),
					huh.NewOption(string(ViewGames), ViewGames),
				).
				Value(&gamesAction).
				Run()

			if err != nil {
				return fmt.Errorf("%w", err)
			}

			switch gamesAction {
			case AddGame:
				if err := gameadd.AddGame(); err != nil {
					fmt.Printf("Error adding game: %v\n", err)
				}
			case EditGame:
				if err := gameedit.EditGame(); err != nil {
					fmt.Printf("Error editing game: %v\n", err)
				}
			case UpdateGames:
				if err := gameupdate.UpdateGames(); err != nil {
					fmt.Printf("Error updating games: %v\n", err)
				}
			case ViewGames:
				if err := gameview.ViewGames(); err != nil {
					fmt.Printf("Error viewing games: %v\n", err)
				}
			}
		}
	}
}
