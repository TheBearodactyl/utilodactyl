// Package update
package update

import (
	"context"
	"fmt"
	"os"
	"utilodactyl/models"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	reviewOwner = "TheBearodactyl"
	reviewRepo = "bearodactyl.dev"
	reviewTag = "v1.0.0"
	reviewFileName = "reviews.json"
)

func UpdateReviews() error {
	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			if models.Cli.Verbose {
				fmt.Println(".env file not found. Falling back to system environment variables.")
			}
		} else {
			return fmt.Errorf("error loading .env file: %w", err)
		}
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("missing GITHUB_TOKEN environment variable")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	release, _, err := client.Repositories.GetReleaseByTag(ctx, reviewOwner, reviewRepo, reviewTag)
	if err != nil {
		return fmt.Errorf("error getting release by tag %s: %w", reviewTag, err)
	}

	for _, asset := range release.Assets {
		if asset.GetName() == reviewFileName {
			if models.Cli.Verbose {
				fmt.Printf("Found existing asset '%s' with ID %d. Deleting...\n", reviewFileName, asset.GetID())
			}
			_, err := client.Repositories.DeleteReleaseAsset(ctx, reviewOwner, reviewRepo, asset.GetID())
			if err != nil {
				return fmt.Errorf("error deleting existing asset %s (ID: %d): %w", reviewFileName, asset.GetID(), err)
			}
			if models.Cli.Verbose {
				fmt.Println("Asset deleted successfully.")
			}
			break
		}
	}

	filePath := "./" + reviewFileName
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open local %s: %w", reviewFileName, err)
	}
	defer file.Close()

	if models.Cli.Verbose {
		fmt.Printf("Uploading new asset '%s' to release ID %d...\n", reviewFileName, release.GetID())
	}
	_, _, err = client.Repositories.UploadReleaseAsset(ctx, reviewOwner, reviewRepo, release.GetID(), &github.UploadOptions{
		Name: reviewFileName,
	}, file)
	if err != nil {
		return fmt.Errorf("error uploading asset: %w", err)
	}

	if models.Cli.Verbose {
		fmt.Printf("Upload of '%s' successful to release ID %d.\n", reviewFileName, release.GetID())
	}

	return nil
}

