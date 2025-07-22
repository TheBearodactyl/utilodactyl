// Package add
package add

import (
	"fmt"
	"strings"
	"utilodactyl/models"
	"utilodactyl/utils"

	"github.com/charmbracelet/huh"
)

func AddBook() error {
	books, err := utils.LoadBooks()
	if err != nil {
		return fmt.Errorf("failed to load books: %v", err)
	}

	var newBook models.Book
	var rating int
	var status string

	basicDetailsGroup := huh.NewGroup(
		huh.NewInput().Title("Title:").Value(&newBook.Title).Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("empty title")
			}

			return nil
		},
		),
		huh.NewInput().
			Title("Author:").
			Value(&newBook.Author).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("author cannot be empty")
				}
				return nil
			}),

		huh.NewConfirm().
			Title("Explicit Content:").
			Value(&newBook.Explicit),

		huh.NewInput().
			Title("Cover Image URL:").
			Value(&newBook.CoverImage).
			Validate(utils.ValidateURL),
		huh.NewText().
			Title("Description:").
			Value(&newBook.Description).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("description cannot be empty")
				}
				return nil
			}),

		huh.NewText().
			Title("Your Thoughts:").
			Value(&newBook.MyThoughts).
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
			Title("Reading Status:").
			Options(
				huh.NewOption("Reading", "Reading"),
				huh.NewOption("Finished", "Finished"),
				huh.NewOption("Plan to Read", "Plan to Read"),
				huh.NewOption("Dropped", "Dropped"),
			).
			Value(&status),
		huh.NewInput().
			Title("Border Color:").
			Value(&newBook.Color).
			Validate(func(s string) error {
				if !utils.IsColorCode(s) {
					return fmt.Errorf("string is not a valid color code")
				}

				return nil
			}),
	)

	if err = huh.NewForm(basicDetailsGroup).Run(); err != nil {
		return fmt.Errorf("form input error for basic book details: %w", err)
	}

	newBook.Rating = uint16(rating)
	newBook.Status = status

	if err = handleGenres(books, &newBook); err != nil {
		return fmt.Errorf("error handling book genres: %w", err)
	}

	if err = handleTags(books, &newBook); err != nil {
		return fmt.Errorf("error handling book tags: %w", err)
	}

	if err = handleLinks(&newBook); err != nil {
		return fmt.Errorf("error handling book links: %w", err)
	}

	newBook.ID, err = utils.GenerateBookID()
	if err != nil {
		return fmt.Errorf("failed to generate unique book ID: %w", err)
	}

	books = append(books, newBook)

	if err = utils.SaveBooks(books); err != nil {
		return fmt.Errorf("failed to save books after adding new entry: %w", err)
	}

	fmt.Println("✅. Book added successfully!")
	return nil
}

func handleGenres(existingBooks []models.Book, book *models.Book) error {
	existingGenres := utils.CollectUniqueBookGenres(existingBooks)
	if len(existingGenres) > 0 {
		err := huh.NewMultiSelect[string]().
			Title("Select existing genres:").
			Options(huh.NewOptions(existingGenres...)...).
			Value(&book.Genres).
			Height(10).
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

func handleTags(existingBooks []models.Book, book *models.Book) error {
	existingTags := utils.CollectUniqueBookTags(existingBooks)
	if len(existingTags) > 0 {
		err := huh.NewMultiSelect[string]().
			Title("Select existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&book.Tags).
			Height(10).
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
			book.Tags = append(book.Tags, customTag)

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

func handleLinks(book *models.Book) error {
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
