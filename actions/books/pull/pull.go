package pull

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	owner    = "TheBearodactyl"
	repo     = "bearodactyl.dev"
	tag      = "v1.0.0"
	fileName = "books.json"
)

func PullBooks() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN not set in .env or environment")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	release, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return fmt.Errorf("failed to get release by tag: %w", err)
	}

	var assetID int64
	for _, asset := range release.Assets {
		if asset.GetName() == fileName {
			assetID = asset.GetID()
			break
		}
	}

	if assetID == 0 {
		return fmt.Errorf("asset %s not found in release %s", fileName, tag)
	}

	rc, url, err := client.Repositories.DownloadReleaseAsset(ctx, owner, repo, assetID)
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

	out, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fileName, err)
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	if err != nil {
		return fmt.Errorf("failed to write asset to file: %w", err)
	}

	fmt.Printf("Downloaded %s successfully\n", fileName)
	return nil
}
