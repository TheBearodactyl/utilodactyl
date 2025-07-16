package view

import (
	"fmt"
	"siteutil/utils"
	"strings"
)

func ViewGames() error {
	games, err := utils.LoadGames()
	if err != nil {
		return fmt.Errorf("failed to load games for viewing: %w", err)
	}

	if len(games) == 0 {
		fmt.Println("No games to show.")
	} else {
		for _, game := range games {
			fmt.Printf("\n📖 %s by %s\n", game.Title, game.Developer)
			fmt.Printf("⭐ Rating: %d\n", game.Rating)
			fmt.Printf("📚 Genres: %s\n", joinStringSlice(game.Genres, ", "))
			fmt.Printf("🏷️ Tags: %s\n", joinStringSlice(game.Tags, ", "))
			fmt.Printf("📄 Description: %s\n", game.Description)
			fmt.Printf("💭 Thoughts: %s\n", game.MyThoughts)
			fmt.Printf("📈 Status: %s\n", game.Status)
			if game.Explicit {
				fmt.Println("🔞 Explicit Content: Yes")
			} else {
				fmt.Println("✅ Explicit Content: No")
			}
			if len(game.Links) > 0 {
				fmt.Println("🔗 Links:")
				for _, link := range game.Links {
					fmt.Printf("  • %s → %s\n", link.Title, link.URL)
				}
			}
		}
	}
	return nil
}

func joinStringSlice(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.Join(s, sep)
}
