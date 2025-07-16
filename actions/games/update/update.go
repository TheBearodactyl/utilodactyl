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
	owner    = "TheBearodactyl"
	repo     = "bearodactyl.dev"
	tag      = "v1.0.0"
	fileName = "games.json"
)

func UpdateGames() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
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

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return fmt.Errorf("error getting release by tag %s: %w", tag, err)
	}

	for _, asset := range release.Assets {
		if asset.GetName() == fileName {
			if models.Cli.Verbose {
				fmt.Printf("Found existing asset '%s' with ID %d. Deleting...\n", fileName, asset.GetID())
			}
			_, err := client.Repositories.DeleteReleaseAsset(ctx, owner, repo, asset.GetID())
			if err != nil {
				return fmt.Errorf("error deleting existing asset %s (ID: %d): %w", fileName, asset.GetID(), err)
			}
			if models.Cli.Verbose {
				fmt.Println("Asset deleted successfully.")
			}
			break
		}
	}

	filePath := "./games.json"
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open local games.json: %w", err)
	}
	defer file.Close()

	if models.Cli.Verbose {
		fmt.Printf("Uploading new asset '%s' to release ID %d...\n", fileName, release.GetID())
	}
	_, _, err = client.Repositories.UploadReleaseAsset(ctx, owner, repo, release.GetID(), &github.UploadOptions{
		Name: fileName,
	}, file)
	if err != nil {
		return fmt.Errorf("error uploading asset: %w", err)
	}

	if models.Cli.Verbose {
		fmt.Printf("Upload of '%s' successful to release ID %d.\n", fileName, release.GetID())
	}

	return nil
}
