package blog

import (
	"fmt"
	"os"
	"sjdhome.com/kurisu-go/database"
)

type Blog struct{}

func (b Blog) getAllPosts() ([]Post, error) {
	var posts []Post
	db := database.GetDatabase()
	rows, err := db.Query("SELECT id, title, created, edited, author, description, commentable, visible, pinToTop FROM posts;")
	if err != nil {
		return nil, fmt.Errorf("Querying 'posts' table: %w", err)
	}
	for rows.Next() {
		var (
			id          string
			title       string
			created     string
			edited      string
			author      string
			description string
			commentable int
			visible     int
			pinToTop    int
		)
		err := rows.Scan(&id, &title, &created, &edited, &author, &description, &commentable, &visible, &pinToTop)
		if err != nil {
			return nil, fmt.Errorf("Extracting data from SQL response: %w", err)
		}
		tags, err := b.getTagsOfPost(id)
		if err != nil {
			return nil, fmt.Errorf("Querying tags of '%s' post: %w", id, err)
		}
		post := Post{
			Id:          id,
			Title:       title,
			Created:     created,
			Edited:      edited,
			Author:      author,
			Tags:        tags,
			Description: description,
			// SQLite 3 has no bool type.
			Commentable: commentable != 0,
			Visible:     visible != 0,
			PinToTop:    pinToTop != 0,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (b Blog) getTagsOfPost(id string) ([]string, error) {
	db := database.GetDatabase()
	tagRows, err := db.Query("SELECT tag FROM tags WHERE post_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("Querying 'post_tag' table where id = '%s': %w", id, err)
	}
	var tags []string
	for tagRows.Next() {
		var tag string
		tagRows.Scan(&tag)
		tags = append(tags, tag)
	}
	return tags, nil
}

func (b Blog) getPost(id string) (*Post, error) {
	db := database.GetDatabase()
	row := db.QueryRow("SELECT title, created, edited, author, description, commentable, visible, pinToTop FROM 'posts' WHERE id = ?", id)
	var (
		title       string
		created     string
		edited      string
		author      string
		description string
		commentable int
		visible     int
		pinToTop    int
	)
	err := row.Scan(&title, &created, &edited, &author, &description, &commentable, &visible, &pinToTop)
	if err != nil {
		return nil, fmt.Errorf("Extracting data from SQL response: %w", err)
	}
	tags, err := b.getTagsOfPost(id)
	if err != nil {
		return nil, fmt.Errorf("Querying tags of '%s' post: %w", id, err)
	}
	post := Post{
		Id:          id,
		Title:       title,
		Created:     created,
		Edited:      edited,
		Author:      author,
		Tags:        tags,
		Description: description,
		// SQLite 3 has no bool type.
		Commentable: commentable != 0,
		Visible:     visible != 0,
		PinToTop:    pinToTop != 0,
	}
	return &post, nil
}

func (b Blog) getPostContent(id string) (string, error) {
	bytes, err := os.ReadFile(fmt.Sprintf("./blog/post/%s.md", id))
	if err != nil {
		return "", fmt.Errorf("Reading blog post '%s' content: %w", id, err)
	}
	return string(bytes), nil
}
