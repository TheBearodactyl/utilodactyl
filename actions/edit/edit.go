package edit

import (
	"bookify/models"
	"bookify/utils"
	"fmt"
	"github.com/charmbracelet/huh"
)

func EditBook() error {
	books, err := utils.LoadBooks()
	if err != nil {
		return err
	}

	if len(books) == 0 {
		fmt.Println("No books available to edit.")
		return nil
	}

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
		return err
	}

	var bookToEdit *models.Book
	for i := range books {
		if books[i].Title == selectedTitle {
			bookToEdit = &books[i]
			break
		}
	}

	if bookToEdit == nil {
		return fmt.Errorf("book not found after selection")
	}

	var rating int = int(bookToEdit.Rating)
	var status string = bookToEdit.Status
	var confirmAddGenre bool
	var customGenre string
	var confirmAddTag bool
	var customTag string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title:").
				Value(&bookToEdit.Title),
			huh.NewInput().
				Title("Author:").
				Value(&bookToEdit.Author),
			huh.NewInput().
				Title("Cover Image:").
				Value(&bookToEdit.CoverImage),
			huh.NewText().
				Title("Description:").
				Value(&bookToEdit.Description),
			huh.NewText().
				Title("Your Thoughts:").
				Value(&bookToEdit.MyThoughts),
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
					huh.NewOption("Unread", "Unread"),
					huh.NewOption("Reading", "Reading"),
					huh.NewOption("Read", "Read"),
				).
				Value(&status),
		),
	).Run()

	if err != nil {
		return err
	}

	bookToEdit.Rating = uint8(rating)
	bookToEdit.Status = status

	// Genres
	existingGenres := utils.CollectUniqueGenres(books)
	if len(existingGenres) > 0 {
		// Use a temporary slice for MultiSelect as it might modify the order or content
		selectedGenres := bookToEdit.Genres
		err = huh.NewMultiSelect[string]().
			Title("Select existing genres:").
			Options(huh.NewOptions(existingGenres...)...).
			Value(&selectedGenres).
			Run()
		if err != nil {
			return err
		}
		bookToEdit.Genres = selectedGenres
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
				Run()
			if err != nil {
				return err
			}
			bookToEdit.Genres = append(bookToEdit.Genres, customGenre)

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
		// Use a temporary slice for MultiSelect
		selectedTags := bookToEdit.Tags
		err = huh.NewMultiSelect[string]().
			Title("Select existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&selectedTags).
			Run()
		if err != nil {
			return err
		}
		bookToEdit.Tags = selectedTags
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
				Run()
			if err != nil {
				return err
			}
			bookToEdit.Tags = append(bookToEdit.Tags, customTag)

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

	err = utils.SaveBooks(books)
	if err != nil {
		return err
	}

	fmt.Println("✅ Book updated.")
	return nil
}
