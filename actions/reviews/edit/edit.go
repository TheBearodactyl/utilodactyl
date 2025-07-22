// Package edit
package edit

import (
	"fmt"
	"strconv"
	"strings"
	"utilodactyl/models"
	"utilodactyl/utils"

	"github.com/charmbracelet/huh"
)

func EditReview() error {
	reviews, err := utils.LoadReviews()
	if err != nil {
		return fmt.Errorf("error loading reviews: %v", err)
	}

	if len(reviews) == 0 {
		return fmt.Errorf("no reviews found to edit")
	}

	reviewOptions := make([]huh.Option[string], len(reviews))
	for i, review := range reviews {
		reviewOptions[i] = huh.NewOption(fmt.Sprintf("Chapter %d: %s", review.Chapter, review.Description), strconv.Itoa(int(review.Chapter)))
	}

	var selectedChapterStr string
	err = huh.NewSelect[string]().
		Title("Choose a review to edit:").
		Options(reviewOptions...).
		Value(&selectedChapterStr).
		Run()
	if err != nil {
		return fmt.Errorf("error selecting review: %v", err)
	}

	selectedChapter, err := strconv.ParseUint(selectedChapterStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid chapter selected: %v", err)
	}

	var reviewToEdit *models.Review
	var originalChapter uint32
	for i := range reviews {
		if reviews[i].Chapter == uint32(selectedChapter) {
			reviewToEdit = &reviews[i]
			originalChapter = reviews[i].Chapter 
			break
		}
	}

	if reviewToEdit == nil {
		return fmt.Errorf("review not found")
	}

	basicDetailsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Chapter Number:").
				Value(new(string)).
				Validate(func(s string) error {
					chapterStr := strings.TrimSpace(s)
					if chapterStr == "" {
						return fmt.Errorf("chapter number cannot be empty")
					}
					chapter, err := parseUint32(chapterStr)
					if err != nil {
						return fmt.Errorf("invalid chapter number: %v", err)
					}
					if chapter == 0 {
						return fmt.Errorf("chapter number must be greater than 0")
					}
					for _, r := range reviews {
						if r.Chapter == chapter && r.Chapter != originalChapter {
							return fmt.Errorf("chapter number %d already exists", chapter)
						}
					}
					reviewToEdit.Chapter = chapter
					return nil
				}),
			huh.NewInput().
				Title("Description:").
				Value(&reviewToEdit.Description).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("description cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Rating (1-5):").
				Value(new(string)).
				Validate(func(s string) error {
					ratingStr := strings.TrimSpace(s)
					if ratingStr == "" {
						return fmt.Errorf("rating cannot be empty")
					}
					rating, err := parseUint8(ratingStr)
					if err != nil {
						return fmt.Errorf("invalid rating: %v", err)
					}
					if rating < 1 || rating > 5 {
						return fmt.Errorf("rating must be between 1 and 5")
					}
					reviewToEdit.Rating = rating
					return nil
				}),
			huh.NewInput().
				Title("Thoughts:").
				Value(&reviewToEdit.Thoughts).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("thoughts cannot be empty")
					}
					return nil
				}),
		),
	)

	if err = basicDetailsForm.Run(); err != nil {
		return fmt.Errorf("error editing review form: %v", err)
	}

	if err = utils.SaveReviews(reviews); err != nil {
		return fmt.Errorf("error saving reviews: %v", err)
	}

	fmt.Printf("✅ Review for Chapter %d updated successfully!\n", reviewToEdit.Chapter)
	return nil
}

func parseUint8(s string) (uint8, error) {
	val, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("failed to parse to number: %w", err)
	}
	return uint8(val), nil
}

func parseUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse to number: %w", err)
	}
	return uint32(val), nil
}
