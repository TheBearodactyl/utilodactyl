// Package view
package view

import (
	"fmt"
	"utilodactyl/utils"
)

func ViewReviews() error {
	reviews, err := utils.LoadReviews()
	if err != nil {
		return fmt.Errorf("failed to load reviews: %w", err)
	}

	if len(reviews) == 0 {
		fmt.Println("No reviews found")
	} else {
		for _, review := range reviews {
			fmt.Printf("\nChapter: %d\n", review.Chapter)
			fmt.Printf("Description: %s\n", review.Description)
			fmt.Printf("Rating: %d/5\n", review.Rating)
			fmt.Printf("Thoughts: %s\n", review.Thoughts)
		}
	}

	return nil
}

