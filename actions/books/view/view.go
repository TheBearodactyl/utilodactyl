package view

import (
	"fmt"
	"siteutil/utils"
	"strings"
)

func ViewBooks() error {
	books, err := utils.LoadBooks()
	if err != nil {
		// Propagate the error from loading books.
		return fmt.Errorf("failed to load books for viewing: %w", err)
	}

	if len(books) == 0 {
		fmt.Println("No books to show.")
	} else {
		// Iterate through each book and print its formatted details.
		for _, book := range books {
			fmt.Printf("\n📖 %s by %s\n", book.Title, book.Author)
			fmt.Printf("⭐ Rating: %d\n", book.Rating)
			fmt.Printf("📚 Genres: %s\n", joinStringSlice(book.Genres, ", "))
			fmt.Printf("🏷️ Tags: %s\n", joinStringSlice(book.Tags, ", "))
			fmt.Printf("📄 Description: %s\n", book.Description)
			fmt.Printf("💭 Thoughts: %s\n", book.MyThoughts)
			fmt.Printf("📈 Status: %s\n", book.Status)
			if book.Explicit {
				fmt.Println("🔞 Explicit Content: Yes")
			} else {
				fmt.Println("✅ Explicit Content: No")
			}
			if len(book.Links) > 0 {
				fmt.Println("🔗 Links:")
				for _, link := range book.Links {
					fmt.Printf("  • %s → %s\n", link.Title, link.URL)
				}
			}
		}
	}
	return nil
}

// joinStringSlice concatenates a slice of strings into a single string,
// with each element separated by the specified separator.
// It returns an empty string if the slice is empty.
func joinStringSlice(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	// Using strings.Join is generally more efficient than manual concatenation in a loop
	// for larger slices, as it pre-calculates the final string size.
	return strings.Join(s, sep)
}
