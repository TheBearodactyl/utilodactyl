// Package utils
package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"utilodactyl/models"
	"sort"
)

const BooksFile = "books.json"
const ProjectsFile = "projects.json"
const GamesFile = "games.json"

func LoadBooks() ([]models.Book, error) {
	data, err := os.ReadFile(BooksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Book{}, nil
		}

		return nil, fmt.Errorf("failed to read %s: %w", BooksFile, err)
	}

	var books []models.Book
	if err = json.Unmarshal(data, &books); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", BooksFile, err)
	}

	return books, nil
}

func LoadGames() ([]models.Game, error) {
	data, err := os.ReadFile(GamesFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Game{}, nil
		}

		return nil, fmt.Errorf("failed to read %s: %w", BooksFile, err)
	}

	var games []models.Game
	if err = json.Unmarshal(data, &games); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", BooksFile, err)
	}

	return games, nil
}

func LoadProjects() ([]models.Project, error) {
	data, err := os.ReadFile(ProjectsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Project{}, nil
		}

		return nil, fmt.Errorf("failed to read %s: %w", ProjectsFile, err)
	}

	var projects []models.Project
	if err = json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", ProjectsFile, err)
	}

	return projects, nil
}

func SaveGames(games []models.Game) error {
	data, err := json.MarshalIndent(games, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal games: %w", err)
	}

	if err = os.WriteFile(GamesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", GamesFile, err)
	}

	return nil
}

func SaveBooks(books []models.Book) error {
	data, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal books: %w", err)
	}

	if err = os.WriteFile(BooksFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", BooksFile, err)
	}

	return nil
}

func SaveProjects(projects []models.Project) error {
	data, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	if err = os.WriteFile(ProjectsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", ProjectsFile, err)
	}

	return nil
}

func GenerateGameID() (uint32, error) {
	games, err := LoadGames()
	if err != nil {
		return 0, fmt.Errorf("failed to load games: %w", err)
	}

	maxID := uint32(0)
	for _, game := range games {
		if game.ID > maxID {
			maxID = game.ID
		}
	}

	return maxID + 1, nil
}

func GenerateBookID() (uint32, error) {
	books, err := LoadBooks()
	if err != nil {
		return 0, fmt.Errorf("failed to load books: %w", err)
	}

	maxID := uint32(0)
	for _, book := range books {
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	return maxID + 1, nil
}

func CollectUniqueBookTags(books []models.Book) []string {
	tagSet := make(map[string]struct{})
	for _, book := range books {
		for _, tag := range book.Tags {
			tagSet[tag] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func CollectUniqueProjectTags(projects []models.Project) []string {
	tagSet := make(map[string]struct{})
	for _, project := range projects {
		for _, tag := range project.Tags {
			tagSet[tag] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func CollectUniqueGameTags(games []models.Game) []string {
	tagSet := make(map[string]struct{})
	for _, game := range games {
		for _, tag := range game.Tags {
			tagSet[tag] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func CollectUniqueBookGenres(books []models.Book) []string {
	genreSet := make(map[string]struct{})
	for _, book := range books {
		for _, genre := range book.Genres {
			genreSet[genre] = struct{}{}
		}
	}

	genres := make([]string, 0, len(genreSet))
	for genre := range genreSet {
		genres = append(genres, genre)
	}
	sort.Strings(genres)
	return genres
}

func CollectUniqueGameGenres(games []models.Game) []string {
	genreSet := make(map[string]struct{})
	for _, game := range games {
		for _, genre := range game.Genres {
			genreSet[genre] = struct{}{}
		}
	}

	genres := make([]string, 0, len(genreSet))
	for genre := range genreSet {
		genres = append(genres, genre)
	}
	sort.Strings(genres)
	return genres
}
