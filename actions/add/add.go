package add

import (
	"bookify/models"
	"bookify/utils"
	"fmt"
	"github.com/charmbracelet/huh"
)

func AddBook() error {
	books, err := utils.LoadBooks()
	if err != nil {
		return err
	}

	var newBook models.Book
	var rating int
	var status string
	var confirmAddGenre bool
	var customGenre string
	var confirmAddTag bool
	var customTag string
	var confirmAddLink bool
	var linkTitle, linkURL string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title:").
				Value(&newBook.Title).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("title cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Author:").
				Value(&newBook.Author).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("author cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("Cover Image URL:").
				Value(&newBook.CoverImage).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("cover image URL cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Title("Description:").
				Value(&newBook.Description).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("description cannot be empty")
					}
					return nil
				}),

			huh.NewText().
				Title("Your Thoughts:").
				Value(&newBook.MyThoughts).
				Validate(func(s string) error {
					if s == "" {
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
		),
	).Run()

	if err != nil {
		return err
	}
	newBook.Rating = uint8(rating)
	newBook.Status = status

	// Genres
	existingGenres := utils.CollectUniqueGenres(books)
	if len(existingGenres) > 0 {
		err = huh.NewMultiSelect[string]().
			Title("Select existing genres:").
			Options(huh.NewOptions(existingGenres...)...).
			Value(&newBook.Genres).
			Run()
		if err != nil {
			return err
		}
	}

	err = huh.NewConfirm().
		Title("Add custom genres?").
		Value(&confirmAddGenre).
		Run()
	if err != nil {
		return err
	}
	if confirmAddGenre {
		for {
			err = huh.NewInput().
				Title("Genre:").
				Value(&customGenre).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("genre cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			newBook.Genres = append(newBook.Genres, customGenre)

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

	// Tags
	existingTags := utils.CollectUniqueTags(books)
	if len(existingTags) > 0 {
		err = huh.NewMultiSelect[string]().
			Title("Select existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&newBook.Tags).
			Run()
		if err != nil {
			return err
		}
	}

	err = huh.NewConfirm().
		Title("Add custom tags?").
		Value(&confirmAddTag).
		Run()
	if err != nil {
		return err
	}
	if confirmAddTag {
		for {
			err = huh.NewInput().
				Title("Tag:").
				Value(&customTag).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("tag cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			newBook.Tags = append(newBook.Tags, customTag)

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

	// Links
	err = huh.NewConfirm().
		Title("Add links?").
		Value(&confirmAddLink).
		Run()
	if err != nil {
		return err
	}
	if confirmAddLink {
		for {
			err = huh.NewInput().
				Title("Link title:").
				Value(&linkTitle).
				Validate(func(s string) error {
					if s == "" {
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
					if s == "" {
						return fmt.Errorf("link URL cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			newBook.Links = append(newBook.Links, models.BookLink{Title: linkTitle, URL: linkURL})

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

	newBook.ID, err = utils.GenerateBookID()
	if err != nil {
		return err
	}

	books = append(books, newBook)
	err = utils.SaveBooks(books)
	if err != nil {
		return err
	}

	fmt.Println("✅ Book added!")
	return nil
}
