package edit

import (
	"fmt"
	"strings"
	"utilodactyl/models"
	"utilodactyl/utils"

	"github.com/charmbracelet/huh"
)

func EditGame() error {
	games, err := utils.LoadGames()
	if err != nil {
		return fmt.Errorf("failed to load games for editing: %w", err)
	}

	if len(games) == 0 {
		fmt.Println("No games available to edit.")
		return nil
	}

	// Prepare options for selecting a book by title.
	gameTitles := make([]string, len(games))
	for i, b := range games {
		gameTitles[i] = b.Title
	}

	var selectedTitle string
	err = huh.NewSelect[string]().
		Title("Choose a game to edit:").
		Options(huh.NewOptions(gameTitles...)...).
		Value(&selectedTitle).
		Run()
	if err != nil {
		return fmt.Errorf("game selection cancelled or failed: %w", err)
	}

	var gameToEdit *models.Game
	for i := range games {
		if games[i].Title == selectedTitle {
			gameToEdit = &games[i]
			break
		}
	}

	if gameToEdit == nil {
		return fmt.Errorf("internal error: selected game '%s' not found", selectedTitle)
	}

	rating := int(gameToEdit.Rating)
	status := gameToEdit.Status

	basicDetailsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title:").
				Value(&gameToEdit.Title).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("title cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Developer:").
				Value(&gameToEdit.Developer).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("author cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Explicit Content:").
				Value(&gameToEdit.Explicit),
			huh.NewInput().
				Title("Cover Image URL:").
				Value(&gameToEdit.CoverImage).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("cover image URL cannot be empty")
					}
					return nil
				}),
			huh.NewText().
				Title("Description:").
				Value(&gameToEdit.Description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("description cannot be empty")
					}
					return nil
				}),
			huh.NewText().
				Title("Your Thoughts:").
				Value(&gameToEdit.MyThoughts).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("your thoughts cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[int]().
				Title("Rating (1–5):").
				Options(
					huh.NewOption("1", 1),
					huh.NewOption("2", 2),
					huh.NewOption("3", 3),
					huh.NewOption("4", 4),
					huh.NewOption("5", 5),
				).
				Value(&rating),
			huh.NewSelect[string]().
				Title("Play Status:").
				Options(
					huh.NewOption("Playing", "Playing"),
					huh.NewOption("Finished", "Finished"),
					huh.NewOption("Plan to Play", "Plan to Play"),
					huh.NewOption("Dropped", "Dropped"),
				).
				Value(&status),
		),
	)

	if err = basicDetailsForm.Run(); err != nil {
		return fmt.Errorf("form input error for game details: %w", err)
	}

	gameToEdit.Rating = uint16(rating)
	gameToEdit.Status = status

	if err = editGenres(games, gameToEdit); err != nil {
		return fmt.Errorf("error editing game genres: %w", err)
	}

	if err = editTags(games, gameToEdit); err != nil {
		return fmt.Errorf("error editing game tags: %w", err)
	}

	if err = editLinks(gameToEdit); err != nil {
		return fmt.Errorf("error editing game links: %w", err)
	}

	if err = utils.SaveGames(games); err != nil {
		return fmt.Errorf("failed to save games after editing: %w", err)
	}

	fmt.Println("✅ Game updated successfully!")
	return nil
}

func editGenres(allGames []models.Game, game *models.Game) error {
	existingGenres := utils.CollectUniqueGameGenres(allGames)
	if len(existingGenres) > 0 {
		selectedGenres := game.Genres
		err := huh.NewMultiSelect[string]().
			Title("Select/Deselect existing genres:").
			Options(huh.NewOptions(existingGenres...)...).
			Value(&selectedGenres).
			Run()
		if err != nil {
			return err
		}
		game.Genres = selectedGenres
	}

	var confirmAddGenre bool
	err := huh.NewConfirm().
		Title("Add custom genres?").
		Value(&confirmAddGenre).
		Run()
	if err != nil {
		return err
	}

	if confirmAddGenre {
		for {
			var customGenre string
			err = huh.NewInput().
				Title("New genre:").
				Value(&customGenre).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("genre cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			game.Genres = append(game.Genres, customGenre)

			var addAnother bool
			err = huh.NewConfirm().
				Title("Add another genre?").
				Value(&addAnother).
				Run()
			if err != nil {
				return err
			}
			if !addAnother {
				break
			}
		}
	}
	return nil
}

func editTags(allGames []models.Game, game *models.Game) error {
	existingTags := utils.CollectUniqueGameTags(allGames)
	if len(existingTags) > 0 {
		selectedTags := game.Tags
		err := huh.NewMultiSelect[string]().
			Title("Select/Deselect existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&selectedTags).
			Run()
		if err != nil {
			return err
		}
		game.Tags = selectedTags
	}

	var confirmAddTag bool
	err := huh.NewConfirm().
		Title("Add custom tags?").
		Value(&confirmAddTag).
		Run()
	if err != nil {
		return err
	}

	if confirmAddTag {
		for {
			var customTag string
			err = huh.NewInput().
				Title("New tag:").
				Value(&customTag).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("tag cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			game.Tags = append(game.Tags, customTag)

			var addAnother bool
			err = huh.NewConfirm().
				Title("Add another tag?").
				Value(&addAnother).
				Run()
			if err != nil {
				return err
			}
			if !addAnother {
				break
			}
		}
	}
	return nil
}

func editLinks(game *models.Game) error {
	var confirmAddLink bool
	err := huh.NewConfirm().
		Title("Add more links?").
		Value(&confirmAddLink).
		Run()
	if err != nil {
		return err
	}

	if confirmAddLink {
		for {
			var linkTitle, linkURL string
			err = huh.NewInput().
				Title("Link title:").
				Value(&linkTitle).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("link title cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			err = huh.NewInput().
				Title("Link URL:").
				Value(&linkURL).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("link URL cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			game.Links = append(game.Links, models.ItemLink{Title: linkTitle, URL: linkURL})

			var addAnother bool
			err = huh.NewConfirm().
				Title("Add another link?").
				Value(&addAnother).
				Run()
			if err != nil {
				return err
			}
			if !addAnother {
				break
			}
		}
	}
	return nil
}
