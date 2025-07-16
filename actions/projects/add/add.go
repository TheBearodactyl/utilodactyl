// Package add
package add

import (
	"fmt"
	"strings"
	"utilodactyl/models"
	"utilodactyl/utils"

	"github.com/charmbracelet/huh"
)

func AddProject() error {
	projects, err := utils.LoadProjects()
	if err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	var newProject models.Project

	basicDetailsGroup := huh.NewGroup(
		huh.NewInput().Title("Name:").Value(&newProject.Name).Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("name is empty")
			}

			return nil
		}),
		huh.NewInput().Title("Description:").Value(&newProject.Description).Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("description is empty")
			}
			return nil
		}),
		huh.NewInput().Title("Source:").Value(&newProject.Source).Validate(utils.ValidateURL),
	)

	if err = huh.NewForm(basicDetailsGroup).Run(); err != nil {
		return fmt.Errorf("error creating new form: %v", err)
	}

	if err = handleTags(projects, &newProject); err != nil {
		return fmt.Errorf("error handling tags: %v", err)
	}

	projects = append(projects, newProject)
	if err = utils.SaveProjects(projects); err != nil {
		return fmt.Errorf("error saving projects: %v", err)
	}

	fmt.Println("Projects saved")
	return nil
}

func handleTags(existingProjects []models.Project, project *models.Project) error {
	existingTags := utils.CollectUniqueProjectTags(existingProjects)
	if len(existingTags) > 0 {
		err := huh.NewMultiSelect[string]().
			Title("Select existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&project.Tags).
			Run()
		if err != nil {
			return err
		}
	}

	var confirmAddTag bool
	err := huh.NewConfirm().
		Title("Add custom tags?").
		Value(&confirmAddTag).
		Run()
	if err != nil {
		return err
	}

	if confirmAddTag {
		for {
			var customTag string
			err = huh.NewInput().
				Title("New tag:").
				Value(&customTag).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("tag cannot be empty")
					}
					return nil
				}).
				Run()
			if err != nil {
				return err
			}
			project.Tags = append(project.Tags, customTag)

			var addAnother bool
			err = huh.NewConfirm().
				Title("Add another tag?").
				Value(&addAnother).
				Run()
			if err != nil {
				return err
			}
			if !addAnother {
				break
			}
		}
	}
	return nil
}
