package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"utilodactyl/models"
)

const (
	booksFile    = "books.json"
	projectsFile = "projects.json"
	gamesFile    = "games.json"
)

func readJSONFile[T any](path string) ([]T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var result []T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", path, err)
	}
	return result, nil
}

func writeJSONFile[T any](path string, items []T) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}

func LoadBooks() ([]models.Book, error) {
	return readJSONFile[models.Book](booksFile)
}

func LoadGames() ([]models.Game, error) {
	return readJSONFile[models.Game](gamesFile)
}

func LoadProjects() ([]models.Project, error) {
	return readJSONFile[models.Project](projectsFile)
}

func SaveBooks(books []models.Book) error {
	return writeJSONFile(booksFile, books)
}

func SaveGames(games []models.Game) error {
	return writeJSONFile(gamesFile, games)
}

func SaveProjects(projects []models.Project) error {
	return writeJSONFile(projectsFile, projects)
}

func generateNextID[T any](loadFunc func() ([]T, error), getID func(T) uint32) (uint32, error) {
	items, err := loadFunc()
	if err != nil {
		return 0, err
	}

	var maxID uint32
	for i := range items {
		if id := getID(items[i]); id > maxID {
			maxID = id
		}
	}
	return maxID + 1, nil
}

func GenerateBookID() (uint32, error) {
	return generateNextID(LoadBooks, func(b models.Book) uint32 { return b.ID })
}

func GenerateGameID() (uint32, error) {
	return generateNextID(LoadGames, func(g models.Game) uint32 { return g.ID })
}

// Common tag/genre collector
func collectUniqueStrings[T any](items []T, extract func(T) []string) []string {
	unique := make(map[string]struct{})
	for _, item := range items {
		for _, s := range extract(item) {
			unique[s] = struct{}{}
		}
	}
	result := make([]string, 0, len(unique))
	for s := range unique {
		result = append(result, s)
	}
	sort.Strings(result)
	return result
}

// Book
func CollectUniqueBookTags(books []models.Book) []string {
	return collectUniqueStrings(books, func(b models.Book) []string { return b.Tags })
}

func CollectUniqueBookGenres(books []models.Book) []string {
	return collectUniqueStrings(books, func(b models.Book) []string { return b.Genres })
}

// Project
func CollectUniqueProjectTags(projects []models.Project) []string {
	return collectUniqueStrings(projects, func(p models.Project) []string { return p.Tags })
}

// Game
func CollectUniqueGameTags(games []models.Game) []string {
	return collectUniqueStrings(games, func(g models.Game) []string { return g.Tags })
}

func CollectUniqueGameGenres(games []models.Game) []string {
	return collectUniqueStrings(games, func(g models.Game) []string { return g.Genres })
}
