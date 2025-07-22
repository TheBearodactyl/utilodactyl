// Package add
package add

import (
	"fmt"
	"strconv"
	"strings"
	"utilodactyl/models"
	"utilodactyl/utils"

	"github.com/charmbracelet/huh"
)

func AddReview() error {
	reviews, err := utils.LoadReviews()
	if err != nil {
		return fmt.Errorf("error loading reviews: %v", err)
	}

	var newReview models.Review

	newReview.Chapter, err = utils.GenerateReviewID()
	if err != nil {
		return fmt.Errorf("error generating review ID: %v", err)
	}

	basicDetailsGroup := huh.NewGroup(
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
					if r.Chapter == chapter {
						return fmt.Errorf("chapter number %d already exists", chapter)
					}
				}
				newReview.Chapter = chapter 
				return nil
			}),
		huh.NewInput().
			Title("Description:").
			Value(&newReview.Description).
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
				newReview.Rating = rating
				return nil
			}),
		huh.NewInput().
			Title("Thoughts:").
			Value(&newReview.Thoughts).
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("thoughts cannot be empty")
				}
				return nil
			}),
	)

	if err = huh.NewForm(basicDetailsGroup).Run(); err != nil {
		return fmt.Errorf("error creating new review form: %v", err)
	}

	reviews = append(reviews, newReview)
	if err = utils.SaveReviews(reviews); err != nil {
		return fmt.Errorf("error saving reviews: %v", err)
	}

	fmt.Printf("✅ Review for Chapter %d saved successfully!\n", newReview.Chapter)
	return nil
}

func parseUint8(s string) (uint8, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0, err
	}
	if i < 0 || i > 255 {
		return 0, fmt.Errorf("value out of uint8 range")
	}
	return uint8(i), nil
}

func parseUint32(s string) (uint32, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse to number: %w", err)
	}
	return uint32(val), nil
}
