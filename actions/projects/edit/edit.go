// Package edit
package edit

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"siteutil/models"
	"siteutil/utils"
	"strings"
)

func EditProject() error {
	projects, err := utils.LoadProjects()
	if err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	if len(projects) == 0 {
		return fmt.Errorf("no projects loaded")
	}

	projectNames := make([]string, len(projects))
	for i, project := range projects {
		projectNames[i] = project.Name
	}

	var selectedName string
	err = huh.NewSelect[string]().
		Title("Choose a project to edit:").
		Options(huh.NewOptions(projectNames...)...).
		Value(&selectedName).
		Run()

	if err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	var projToEdit *models.Project
	for i := range projects {
		if projects[i].Name == selectedName {
			projToEdit = &projects[i]
			break
		}
	}

	if projToEdit == nil {
		return fmt.Errorf("project not found")
	}

	basicDetailsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name:").Value(&projToEdit.Name).Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("name cannot be empty")
				}
				return nil
			}),
			huh.NewInput().Title("Description:").Value(&projToEdit.Description).Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("description cannot be empty")
				}
				return nil
			}),
			huh.NewInput().Title("Source Repo:").Value(&projToEdit.Source).Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("source repo cannot be empty")
				}
				return nil
			}),
		),
	)

	if err = basicDetailsForm.Run(); err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	if err = editTags(projects, projToEdit); err != nil {
		return fmt.Errorf("error loading projects: %v", err)
	}

	if err = utils.SaveProjects(projects); err != nil {
		return fmt.Errorf("error saving projects: %v", err)
	}

	fmt.Println("✅ Project updated successfully!")
	return nil
}

func editTags(allBooks []models.Project, project *models.Project) error {
	existingTags := utils.CollectUniqueProjectTags(allBooks)
	if len(existingTags) > 0 {
		selectedTags := project.Tags
		err := huh.NewMultiSelect[string]().
			Title("Select/Deselect existing tags:").
			Options(huh.NewOptions(existingTags...)...).
			Value(&selectedTags).
			Run()
		if err != nil {
			return err
		}
		project.Tags = selectedTags
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
