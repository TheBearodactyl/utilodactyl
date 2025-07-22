// Package pull
package pull

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"utilodactyl/models"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	reviewOwner    = "TheBearodactyl"
	reviewRepo     = "bearodactyl.dev"
	reviewTag      = "v1.0.0"
	reviewFileName = "reviews.json"
)

func PullReviews() error {
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
		return fmt.Errorf("failed to get release by tag %s: %w", reviewTag, err)
	}

	var assetID int64
	for _, asset := range release.Assets {
		if asset.GetName() == reviewFileName {
			assetID = asset.GetID()
			break
		}
	}

	if assetID == 0 {
		return fmt.Errorf("asset %s not found in release %s", reviewFileName, reviewTag)
	}

	// Download the asset.
	rc, url, err := client.Repositories.DownloadReleaseAsset(ctx, reviewOwner, reviewRepo, assetID)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer func() {
		if rc != nil {
			rc.Close()
		}
	}()

	var data io.ReadCloser
	if rc != nil {
		data = rc
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch asset from redirect URL: %w", err)
		}
		data = resp.Body
		defer resp.Body.Close()
	}

	out, err := os.Create(reviewFileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", reviewFileName, err)
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	if err != nil {
		return fmt.Errorf("failed to write asset to file: %w", err)
	}

	fmt.Printf("Downloaded %s successfully\n", reviewFileName)
	return nil
}
