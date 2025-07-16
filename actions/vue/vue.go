package vue

import (
	"bookify/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateVueData() error {
	books, err := utils.LoadBooks()
	if err != nil {
		return err
	}

	if len(books) == 0 {
		fmt.Println("No books to export.")
		return nil
	}

	var output strings.Builder
	output.WriteString("export const books = ref([\n")

	for _, book := range books {
		output.WriteString("    {\n")
		output.WriteString(fmt.Sprintf("        id: %d,\n", book.ID))
		output.WriteString(fmt.Sprintf("        title: \"%s\",\n", escapeString(book.Title)))
		output.WriteString(fmt.Sprintf("        author: \"%s\",\n", escapeString(book.Author)))
		output.WriteString("        genres: [")
		genres := make([]string, len(book.Genres))
		for i, g := range book.Genres {
			genres[i] = fmt.Sprintf("\"%s\"", escapeString(g))
		}
		output.WriteString(strings.Join(genres, ", "))
		output.WriteString("],\n")
		output.WriteString(fmt.Sprintf("        rating: %d,\n", book.Rating))
		output.WriteString(fmt.Sprintf("        status: \"%s\",\n", escapeString(book.Status)))
		output.WriteString(fmt.Sprintf("        coverImage: \"%s\",\n", escapeString(book.CoverImage)))
		output.WriteString(fmt.Sprintf("        description: `%s`,\n", escapeString(book.Description)))
		output.WriteString(fmt.Sprintf("        myThoughts: `%s`,\n", escapeString(book.MyThoughts)))
		output.WriteString("        tags: [")
		tags := make([]string, len(book.Tags))
		for i, t := range book.Tags {
			tags[i] = fmt.Sprintf("\"%s\"", escapeString(t))
		}
		output.WriteString(strings.Join(tags, ", "))
		output.WriteString("],\n")
		output.WriteString("        links: [\n")
		for _, link := range book.Links {
			output.WriteString(fmt.Sprintf(
				"            { title: \"%s\", url: \"%s\" },\n",
				escapeString(link.Title),
				escapeString(link.URL),
			))
		}
		output.WriteString("        ],\n")
		output.WriteString("    },\n")
	}

	output.WriteString("])\n")

	// Create 'data' directory if it doesn't exist
	outputDir := "data"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", outputDir, err)
		}
	}

	filePath := filepath.Join(outputDir, "books_data.js")
	err = os.WriteFile(filePath, []byte(output.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Vue data to file: %w", err)
	}

	fmt.Println("✅ Vue ref data written to books_data.js")
	return nil
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "`", "\\`")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
