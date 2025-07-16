package models

type Book struct {
	ID          uint32     `json:"id"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Genres      []string   `json:"genres"`
	Rating      uint8      `json:"rating"`
	CoverImage  string     `json:"coverImage"`
	Description string     `json:"description"`
	MyThoughts  string     `json:"myThoughts"`
	Tags        []string   `json:"tags"`
	Links       []BookLink `json:"links"`
	Status      string     `json:"status"`
}

type BookLink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
