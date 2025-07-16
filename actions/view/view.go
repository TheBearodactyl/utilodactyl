package view

import (
	"bookify/utils"
	"fmt"
)

func ViewBooks() error {
	books, err := utils.LoadBooks()
	if err != nil {
		return err
	}

	if len(books) == 0 {
		fmt.Println("No books to show.")
	} else {
		for _, book := range books {
			fmt.Printf("\n📖 %s by %s\n", book.Title, book.Author)
			fmt.Printf("⭐ Rating: %d\n", book.Rating)
			fmt.Printf("📚 Genres: %s\n", joinStringSlice(book.Genres, ", "))
			fmt.Printf("🏷️ Tags: %s\n", joinStringSlice(book.Tags, ", "))
			fmt.Printf("📄 Description: %s\n", book.Description)
			fmt.Printf("💭 Thoughts: %s\n", book.MyThoughts)
			fmt.Printf("📈 Status: %s\n", book.Status)
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

func joinStringSlice(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	result := s[0]
	for i := 1; i < len(s); i++ {
		result += sep + s[i]
	}
	return result
}
