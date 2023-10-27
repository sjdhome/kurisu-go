package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"regexp"
)

type Blog struct {
	db *sql.DB
}

type PostMetadata struct {
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

func (b Blog) GetAllPosts() ([]PostMetadata, error) {
	var posts []PostMetadata
	rows, err := b.db.Query("SELECT id, title, created, edited, author, tags, description, commentable, visible, pinToTop FROM 'posts';")
	if err != nil {
		return nil, fmt.Errorf("Querying databse: %w", err)
	}
	for rows.Next() {
		var (
			id          string
			title       string
			created     string
			edited      string
			author      string
			tags        []string
			description string
			commentable int
			visible     int
			pinToTop    int
		)
		if err := rows.Scan(&id, &title, &created, &edited, &author, &description, &commentable, &visible, &pinToTop); err != nil {
			return nil, fmt.Errorf("Extracting data from result: %w", err)
		}
		tagRows, err := b.db.Query("SELECT tag FROM tags WHERE id = ?", id)
		if err != nil {
			return nil, fmt.Errorf("Querying tags for blog post %s: %w", id, err)
		}
		for tagRows.Next() {
			var tag string
			tagRows.Scan(&tag)
			tags = append(tags, tag)
		}
		post := PostMetadata{
			Id:          id,
			Title:       title,
			Created:     created,
			Edited:      edited,
			Author:      author,
			Tags:        tags,
			Description: description,
			Commentable: commentable != 0,
			Visible:     visible != 0,
			PinToTop:    pinToTop != 0,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (b Blog) GetPost(id string) (*PostMetadata, error) {
	row := b.db.QueryRow("SELECT title, created, edited, author, description, commentable, visible, pinToTop FROM 'posts' WHERE id = ?", id)
	var (
		title       string
		created     string
		edited      string
		author      string
		tags        []string
		description string
		commentable int
		visible     int
		pinToTop    int
	)
	if err := row.Scan(&title, &created, &edited, &author, &description, &commentable, &visible, &pinToTop); err != nil {
		return nil, fmt.Errorf("Extracting data from result: %w", err)
	}
	tagRows, err := b.db.Query("SELECT tag FROM tags WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("Querying tags for blog post %s: %w", id, err)
	}
	for tagRows.Next() {
		var tag string
		tagRows.Scan(&tag)
		tags = append(tags, tag)
	}
	post := PostMetadata{
		Id:          id,
		Title:       title,
		Created:     created,
		Edited:      edited,
		Author:      author,
		Tags:        tags,
		Description: description,
		Commentable: commentable != 0,
		Visible:     visible != 0,
		PinToTop:    pinToTop != 0,
	}
	return &post, nil
}

func (b Blog) GetPostContent(id string) (string, error) {
	return "Example", nil
}

type BlogHandler struct {
	blog *Blog
}

func (h BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	const REGEXP_ERROR = "Unable to parse regular expression for URL."
	onError := func(reason string, err error) {
		log.Println(fmt.Errorf("%s: %w.", reason, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	if r.Method == "GET" {
		isSinglePost, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z\-]+\/`, r.URL.Path)
		if err != nil {
			onError(REGEXP_ERROR, err)
			return
		}
		if isSinglePost {
			readContent, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z\-]+\/content\/$`, r.URL.Path)
			if err != nil {
				onError(REGEXP_ERROR, err)
				return
			}
			readMetadata, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z\-]+\/$`, r.URL.Path)
			if err != nil {
				onError(REGEXP_ERROR, err)
				return
			}
			if readContent {
				// Read single post content
				re := regexp.MustCompile(`\/blog\/post\/([A-Za-z\-]+)\/content\/$`)
				id := string(re.FindSubmatch([]byte(r.URL.Path))[1])
				postContent, err := h.blog.GetPostContent(id)
				if err != nil {
					onError(fmt.Sprintf("Unable to read blog post %s content", id), err)
					return
				}
				w.Write([]byte(postContent))
			} else if readMetadata {
				// Get single post
				re := regexp.MustCompile(`\/blog\/post\/([A-Za-z\-]+)\/$`)
				id := string(re.FindSubmatch([]byte(r.URL.Path))[1])
				post, err := h.blog.GetPost(id)
				if err != nil {
					onError(fmt.Sprintf("Unable to get blog post %s", id), err)
					return
				}
				j, err := json.Marshal(post)
				if err != nil {
					onError(fmt.Sprintf("Unable to marshal blog post %s", id), err)
					return
				}
				w.Write(j)
			}
		} else if r.URL.Path == "/blog/post/" {
			// Get all posts
			posts, err := h.blog.GetAllPosts()
			if err != nil {
				onError("Unable to get all blog posts", err)
				return
			}
			j, err := json.Marshal(posts)
			if err != nil {
				onError("Unable to marshal blog posts", err)
				return
			}
			w.Write(j)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		// Not supported
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
