// Package actions
package actions

import (
	"fmt"
	bookadd "utilodactyl/actions/books/add"
	bookedit "utilodactyl/actions/books/edit"
	bookpull "utilodactyl/actions/books/pull"
	bookupdate "utilodactyl/actions/books/update"
	bookview "utilodactyl/actions/books/view"
	gameadd "utilodactyl/actions/games/add"
	gameedit "utilodactyl/actions/games/edit"
	gamepull "utilodactyl/actions/games/pull"
	gameupdate "utilodactyl/actions/games/update"
	gameview "utilodactyl/actions/games/view"
	projectadd "utilodactyl/actions/projects/add"
	projectedit "utilodactyl/actions/projects/edit"
	projectpull "utilodactyl/actions/projects/pull"
	projectupdate "utilodactyl/actions/projects/update"
	projectview "utilodactyl/actions/projects/view"
	reviewadd "utilodactyl/actions/reviews/add"
	reviewedit "utilodactyl/actions/reviews/edit"
	reviewpull "utilodactyl/actions/reviews/pull"
	reviewupdate "utilodactyl/actions/reviews/update"
	reviewview "utilodactyl/actions/reviews/view"

	"github.com/charmbracelet/huh"
)

type AppAction string

const (
	PullAll        AppAction = "Get the latest releases"
	Books          AppAction = "Operate on `books.json`"
	Projects       AppAction = "Operate on `projects.json`"
	Games          AppAction = "Operate on `games.json`"
	Reviews        AppAction = "Operate on `reviews.json`"
	AddBook        AppAction = "Add a new book"
	AddProject     AppAction = "Add a new project"
	AddGame        AppAction = "Add a new game"
	AddReview      AppAction = "Add a new Chapter Review"
	ViewBooks      AppAction = "View existing books"
	ViewProject    AppAction = "View existing projects"
	ViewGames      AppAction = "View existing games"
	ViewReviews    AppAction = "View existing reviews"
	EditBook       AppAction = "Edit a book"
	EditProject    AppAction = "Edit a project"
	EditGame       AppAction = "Edit a game"
	EditReview     AppAction = "Edit a review"
	UpdateBooks    AppAction = "Update the `books.json` release"
	UpdateProjects AppAction = "Update the `projects.json` release"
	UpdateGames    AppAction = "Update the `games.json` release"
	UpdateReviews  AppAction = "Update the `reviews.json` release"
	PullBooks      AppAction = "Pull the latest `books.json` release"
	PullGames      AppAction = "Pull the latest `games.json` release"
	PullProjects   AppAction = "Pull the latest `projects.json` release"
	PullReviews    AppAction = "Pull the latest `reviews.json` release"
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
				huh.NewOption(string(Reviews), Reviews),
				huh.NewOption(string(PullAll), PullAll),
			).
			Value(&action).Run()
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		switch action {
		case PullAll:
			if err := bookpull.PullBooks(); err != nil {
				fmt.Printf("Error pulling books.json: %v\n", err)
			}
			if err := projectpull.PullProjects(); err != nil {
				fmt.Printf("Error pulling projects.json: %v\n", err)
			}
			if err := gamepull.PullGames(); err != nil {
				fmt.Printf("Error pulling games.json: %v\n", err)
			}
			if err := reviewpull.PullReviews(); err != nil {
				fmt.Printf("Error pulling reviews.json: %v\n", err)
			}
		case Books:
			var bookAction AppAction
			err := huh.NewSelect[AppAction]().
				Title("What would you like to do?").
				Options(
					huh.NewOption(string(AddBook), AddBook),
					huh.NewOption(string(ViewBooks), ViewBooks),
					huh.NewOption(string(EditBook), EditBook),
					huh.NewOption(string(PullBooks), PullBooks),
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
			case PullBooks:
				if err := bookpull.PullBooks(); err != nil {
					fmt.Printf("Error pulling books.json: %v\n", err)
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
					huh.NewOption(string(PullProjects), PullProjects),
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
			case PullProjects:
				if err := projectpull.PullProjects(); err != nil {
					fmt.Printf("Error pulling projects.json: %v\n", err)
				}
			}
		case Games:
			var gamesAction AppAction
			err := huh.NewSelect[AppAction]().
				Title("What would you like to do?").
				Options(
					huh.NewOption(string(AddGame), AddGame),
					huh.NewOption(string(EditGame), EditGame),
					huh.NewOption(string(PullGames), PullGames),
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
			case PullGames:
				if err := gamepull.PullGames(); err != nil {
					fmt.Printf("Error pulling games.json: %v\n", err)
				}
			}
		case Reviews:
			var reviewAction AppAction
			err := huh.NewSelect[AppAction]().
				Title("What would you like to do?").
				Options(
					huh.NewOption(string(AddReview), AddReview),
					huh.NewOption(string(EditReview), EditReview),
					huh.NewOption(string(PullReviews), PullReviews),
					huh.NewOption(string(UpdateReviews), UpdateReviews),
					huh.NewOption(string(ViewReviews), ViewReviews),
				).
				Value(&reviewAction).
				Run()
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			switch reviewAction {
			case AddReview:
				if err := reviewadd.AddReview(); err != nil {
					fmt.Printf("Error adding review: %v\n", err)
				}
			case EditReview:
				if err := reviewedit.EditReview(); err != nil {
					fmt.Printf("Error editing review: %v\n", err)
				}
			case UpdateReviews:
				if err := reviewupdate.UpdateReviews(); err != nil {
					fmt.Printf("Error updating review: %v\n", err)
				}
			case ViewReviews:
				if err := reviewview.ViewReviews(); err != nil {
					fmt.Printf("Error viewing review: %v\n", err)
				}
			case PullReviews:
				if err := reviewpull.PullReviews(); err != nil {
					fmt.Printf("Error pulling reviews.json: %v\n", err)
				}
			}

		}

	}
}
