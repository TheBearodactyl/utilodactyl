// woah
package utils

import (
	"bookify/models"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const BooksFileName = "books.json"

func LoadBooks() ([]models.Book, error) {
	data, err := os.ReadFile(BooksFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Book{}, nil // Return empty slice if file does not exist
		}
		return nil, fmt.Errorf("failed to read books.json: %w", err)
	}

	var books []models.Book
	err = json.Unmarshal(data, &books)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal books.json: %w", err)
	}
	return books, nil
}

func SaveBooks(books []models.Book) error {
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal books to JSON: %w", err)
	}

	err = os.WriteFile(BooksFileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write books.json: %w", err)
	}
	return nil
}

func GenerateBookID() (uint32, error) {
	books, err := LoadBooks()
	if err != nil {
		return 0, err
	}

	var maxID uint32 = 0
	for _, book := range books {
		if book.ID > maxID {
			maxID = book.ID
		}
	}
	return maxID + 1, nil
}

func CollectUniqueTags(books []models.Book) []string {
	tagMap := make(map[string]struct{})
	for _, book := range books {
		for _, tag := range book.Tags {
			tagMap[tag] = struct{}{}
		}
	}

	uniqueTags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		uniqueTags = append(uniqueTags, tag)
	}
	sort.Strings(uniqueTags)
	return uniqueTags
}

func CollectUniqueGenres(books []models.Book) []string {
	genreMap := make(map[string]struct{})
	for _, book := range books {
		for _, genre := range book.Genres {
			genreMap[genre] = struct{}{}
		}
	}

	uniqueGenres := make([]string, 0, len(genreMap))
	for genre := range genreMap {
		uniqueGenres = append(uniqueGenres, genre)
	}
	sort.Strings(uniqueGenres)
	return uniqueGenres
}
