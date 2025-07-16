package view

import (
	"fmt"
	"siteutil/utils"
	"strings"
)

func ViewProjects() error {
	projects, err := utils.LoadProjects()
	if err != nil {
		return fmt.Errorf("failed to load projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("No projects found")
	} else {
		for _, project := range projects {
			fmt.Printf("\nName: %s\n", project.Name)
			fmt.Printf("\nDescription: %s\n", project.Description)
			fmt.Printf("\nTags: %v\n", joinStringSlice(project.Tags, ", "))
			fmt.Printf("\nSource Repo: %s\n", project.Source)
		}
	}

	return nil
}

// joinStringSlice concatenates a slice of strings into a single string,
// with each element separated by the specified separator.
// It returns an empty string if the slice is empty.
func joinStringSlice(s []string, sep string) string {
	if len(s) == 0 {
		return ""
	}
	// Using strings.Join is generally more efficient than manual concatenation in a loop
	// for larger slices, as it pre-calculates the final string size.
	return strings.Join(s, sep)
}
