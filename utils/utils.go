// Package utils
package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"utilodactyl/models"

	"github.com/charmbracelet/huh"
)

const (
	booksFile    = "books.json"
	projectsFile = "projects.json"
	gamesFile    = "games.json"
	reviewsFile  = "reviews.json"
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

func LoadReviews() ([]models.Review, error) {
	return readJSONFile[models.Review](reviewsFile)
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

func SaveReviews(reviews []models.Review) error {
	return writeJSONFile(reviewsFile, reviews)
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

func GenerateReviewID() (uint32, error) {
	return generateNextID(LoadReviews, func(r models.Review) uint32 { return r.Chapter })
}

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

func CollectUniqueBookTags(books []models.Book) []string {
	return collectUniqueStrings(books, func(b models.Book) []string { return b.Tags })
}

func CollectUniqueBookGenres(books []models.Book) []string {
	return collectUniqueStrings(books, func(b models.Book) []string { return b.Genres })
}

func CollectUniqueProjectTags(projects []models.Project) []string {
	return collectUniqueStrings(projects, func(p models.Project) []string { return p.Tags })
}

func CollectUniqueGameTags(games []models.Game) []string {
	return collectUniqueStrings(games, func(g models.Game) []string { return g.Tags })
}

func CollectUniqueGameGenres(games []models.Game) []string {
	return collectUniqueStrings(games, func(g models.Game) []string { return g.Genres })
}

func ValidateURL(input string) error {
	input = strings.TrimSpace(input)

	if input == "" {
		return fmt.Errorf("url cannot be empty")
	}

	u, err := url.ParseRequestURI(input)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("invalid url format")
	}

	return nil
}

func IsColorCode(str string) bool {
	match, err := regexp.MatchString("^#?[0-9A-Fa-f]{6}$", str)
	if err != nil {
		return false
	}

	return match
}

func ValidateColor(s string) error {
	if !IsColorCode(s) {
		return fmt.Errorf("string is not a valid hex code")
	}

	return nil
}

func GenPercentOpts() []huh.Option[uint32] {
	options := make([]huh.Option[uint32], 0, 100)

	for i := 1; i <= 100; i++ {
		option := huh.NewOption(fmt.Sprintf("%d", i), uint32(i))
		options = append(options, option)
	}

	return options
}
