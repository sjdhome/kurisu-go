package blog

type Post struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Created     string   `json:"created"`
	Edited      string   `json:"edited"`
	Author      string   `json:"author"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Commentable bool     `json:"commentable"`
	Visible     bool     `json:"visible"`
	PinToTop    bool     `json:"pinToTop"`
}
