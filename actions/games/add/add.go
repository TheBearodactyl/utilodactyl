package add

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"siteutil/models"
	"siteutil/utils"
	"strings"
)

func AddGame() error {
	books, err := utils.LoadGames()
	if err != nil {
		return fmt.Errorf("failed to load books: %v", err)
	}

	var newGame models.Game
	var rating int
	var status string

	basicDetailsGroup := huh.NewGroup(
		huh.NewInput().Title("Title:").Value(&newGame.Title).Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("empty title")
			}

			return nil
		}),
		huh.NewInput().
			Title("Developer:").
			Value(&newGame.Developer).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("developer cannot be empty")
				}
				return nil
			}),

		huh.NewConfirm().
			Title("Explicit Content:").
			Value(&newGame.Explicit),

		huh.NewInput().
			Title("Cover Image URL:").
			Value(&newGame.CoverImage).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("cover image URL cannot be empty")
				}
				return nil
			}),
		huh.NewText().
			Title("Description:").
			Value(&newGame.Description).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("description cannot be empty")
				}
				return nil
			}),

		huh.NewText().
			Title("Your Thoughts:").
			Value(&newGame.MyThoughts).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("your thoughts cannot be empty")
				}
				return nil
			}),

		huh.NewSelect[int]().
			Title("Rating (1-5):").
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
	)

	if err = huh.NewForm(basicDetailsGroup).Run(); err != nil {
		return fmt.Errorf("form input error for basic game details: %w", err)
	}

	newGame.Rating = uint8(rating)
	newGame.Status = status

	if err = handleGenres(books, &newGame); err != nil {
		return fmt.Errorf("error handling book genres: %w", err)
	}

	if err = handleTags(books, &newGame); err != nil {
		return fmt.Errorf("error handling book tags: %w", err)
	}

	if err = handleLinks(&newGame); err != nil {
		return fmt.Errorf("error handling book links: %w", err)
	}

	newGame.ID, err = utils.GenerateGameID()
	if err != nil {
		return fmt.Errorf("failed to generate unique game ID: %w", err)
	}

	books = append(books, newGame)

	if err = utils.SaveGames(books); err != nil {
		return fmt.Errorf("failed to save games after adding new entry: %w", err)
	}

	fmt.Println("✅. Game added successfully!")
	return nil
}

func handleGenres(existingBooks []models.Game, book *models.Game) error {
	existingGenres := utils.CollectUniqueGameGenres(existingBooks)
	if len(existingGenres) > 0 {
		// Allow selecting multiple existing genres.
		err := huh.NewMultiSelect[string]().
			Title("Select existing genres:").
			Options(huh.NewOptions(existingGenres...)...).
			Value(&book.Genres).
			Run()
		if err != nil {
			return err
		}
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
			book.Genres = append(book.Genres, customGenre)

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

func handleTags(existingGames []models.Game, game *models.Game) error {
	existingTags := utils.CollectUniqueGameTags(existingGames)
	if len(existingTags) > 0 {
		err := huh.NewMultiSelect[string]().
			Title("Select existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&game.Tags).
			Run()
		if err != nil {
			return err
		}
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

func handleLinks(book *models.Game) error {
	var confirmAddLink bool
	err := huh.NewConfirm().
		Title("Add links?").
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
			book.Links = append(book.Links, models.ItemLink{Title: linkTitle, URL: linkURL})

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
