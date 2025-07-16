package edit

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"utilodactyl/models"
	"utilodactyl/utils"
	"strings"
)

// EditBook allows the user to select an existing book and modify its details.
// It loads all books, prompts the user to choose a book, displays its current
// details in a form, and saves the updated book list.
func EditBook() error {
	books, err := utils.LoadBooks()
	if err != nil {
		return fmt.Errorf("failed to load books for editing: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("No books available to edit.")
		return nil
	}

	// Prepare options for selecting a book by title.
	bookTitles := make([]string, len(books))
	for i, b := range books {
		bookTitles[i] = b.Title
	}

	var selectedTitle string
	err = huh.NewSelect[string]().
		Title("Choose a book to edit:").
		Options(huh.NewOptions(bookTitles...)...).
		Value(&selectedTitle).
		Run()
	if err != nil {
		return fmt.Errorf("book selection cancelled or failed: %w", err)
	}

	// Find the selected book. We use a pointer to modify the book directly in the slice.
	var bookToEdit *models.Book
	for i := range books {
		if books[i].Title == selectedTitle {
			bookToEdit = &books[i]
			break
		}
	}

	if bookToEdit == nil {
		return fmt.Errorf("internal error: selected book '%s' not found", selectedTitle)
	}

	// Temporary variables to hold form input values.
	rating := int(bookToEdit.Rating)
	status := bookToEdit.Status

	// Define the form for editing basic book details.
	basicDetailsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title:").
				Value(&bookToEdit.Title).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("title cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Author:").
				Value(&bookToEdit.Author).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("author cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Explicit Content:").
				Value(&bookToEdit.Explicit),
			huh.NewInput().
				Title("Cover Image URL:").
				Value(&bookToEdit.CoverImage).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("cover image URL cannot be empty")
					}
					return nil
				}),
			huh.NewText().
				Title("Description:").
				Value(&bookToEdit.Description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("description cannot be empty")
					}
					return nil
				}),
			huh.NewText().
				Title("Your Thoughts:").
				Value(&bookToEdit.MyThoughts).
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
				Title("Reading Status:").
				Options(
					huh.NewOption("Reading", "Reading"),
					huh.NewOption("Finished", "Finished"),
					huh.NewOption("Plan to Read", "Plan to Read"),
					huh.NewOption("Dropped", "Dropped"),
				).
				Value(&status),
		),
	)

	// Run the basic details form.
	if err = basicDetailsForm.Run(); err != nil {
		return fmt.Errorf("form input error for book details: %w", err)
	}

	bookToEdit.Rating = uint8(rating) // Update the book's rating.
	bookToEdit.Status = status        // Update the book's status.

	// Handle genre modifications.
	if err = editGenres(books, bookToEdit); err != nil {
		return fmt.Errorf("error editing book genres: %w", err)
	}

	// Handle tag modifications.
	if err = editTags(books, bookToEdit); err != nil {
		return fmt.Errorf("error editing book tags: %w", err)
	}

	// Handle link modifications (allowing add/remove or edit existing).
	// For simplicity, current implementation re-prompts to add new links.
	// A more robust solution might allow selecting and editing/deleting existing links.
	if err = editLinks(bookToEdit); err != nil {
		return fmt.Errorf("error editing book links: %w", err)
	}

	// Save the updated list of books back to storage.
	if err = utils.SaveBooks(books); err != nil {
		return fmt.Errorf("failed to save books after editing: %w", err)
	}

	fmt.Println("✅ Book updated successfully!")
	return nil
}

// editGenres allows the user to select/deselect existing genres and add new custom ones.
func editGenres(allBooks []models.Book, book *models.Book) error {
	existingGenres := utils.CollectUniqueBookGenres(allBooks)
	if len(existingGenres) > 0 {
		// Pre-select current genres for the book in the multi-select.
		selectedGenres := book.Genres
		err := huh.NewMultiSelect[string]().
			Title("Select/Deselect existing genres:").
			Options(huh.NewOptions(existingGenres...)...).
			Value(&selectedGenres).
			Run()
		if err != nil {
			return err
		}
		book.Genres = selectedGenres // Update the book's genres.
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

// editTags allows the user to select/deselect existing tags and add new custom ones.
func editTags(allBooks []models.Book, book *models.Book) error {
	existingTags := utils.CollectUniqueBookTags(allBooks)
	if len(existingTags) > 0 {
		// Pre-select current tags for the book in the multi-select.
		selectedTags := book.Tags
		err := huh.NewMultiSelect[string]().
			Title("Select/Deselect existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&selectedTags).
			Run()
		if err != nil {
			return err
		}
		book.Tags = selectedTags // Update the book's tags.
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

// editLinks allows adding new links to a book.
// Note: This current implementation only adds new links. To allow full editing (modify/delete),
// a more complex form flow would be needed (e.g., display existing links and allow selecting to edit/delete).
func editLinks(book *models.Book) error {
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
