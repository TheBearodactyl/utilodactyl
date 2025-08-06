// Package models
package models

type Book struct {
	ID          uint32     `json:"id"`          // Unique identifier for the book.
	Title       string     `json:"title"`       // The title of the book.
	Author      string     `json:"author"`      // The author of the book.
	Genres      []string   `json:"genres"`      // A list of genres the book belongs to.
	Rating      uint16     `json:"rating"`      // User rating for the book (1-5).
	CoverImage  string     `json:"coverImage"`  // URL or path to the book's cover image.
	Description string     `json:"description"` // A brief description of the book.
	MyThoughts  string     `json:"myThoughts"`  // User's personal thoughts or review on the book.
	Tags        []string   `json:"tags"`        // A list of tags associated with the book.
	Links       []ItemLink `json:"links"`       // Relevant links for the book (e.g., purchase, review).
	Status      string     `json:"status"`      // Current reading status (e.g., "Reading", "Finished").
	Explicit    bool       `json:"explicit"`    // Indicates if the book contains explicit content.
	Color       string     `json:"color"`
}

type Game struct {
	ID          uint32     `json:"id"`
	Title       string     `json:"title"`
	Developer   string     `json:"developer"`
	Genres      []string   `json:"genres"`
	Tags        []string   `json:"tags"`
	Rating      uint32     `json:"rating"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	MyThoughts  string     `json:"myThoughts"`
	Links       []ItemLink `json:"links"`
	Explicit    bool       `json:"explicit"`
	CoverImage  string     `json:"coverImage"`
	Percent     uint32     `json:"percent"`
}

type Project struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Tags           []string `json:"tags"`
	Source         string   `json:"source"`
	InstallCommand string   `json:"installCommand"`
}

type Review struct {
	Chapter     uint32 `json:"chapter"`
	Description string `json:"description"`
	Rating      uint8  `json:"rating"`
	Thoughts    string `json:"thoughts"`
}

type ItemLink struct {
	Title string `json:"title"` // The title or description of the link.
	URL   string `json:"url"`   // The URL of the link.
}

var Cli struct {
	Verbose bool `arg:"-v,--verbose" help:"Show advanced logs when updating data"`
}
