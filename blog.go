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
	ShowComment bool     `json:"showComment"`
	Visible     bool     `json:"visible"`
	PinToTop    bool     `json:"pinToTop"`
}

func (b Blog) GetAllPosts() ([]PostMetadata, error) {
	var posts []PostMetadata
	rows, err := b.db.Query("SELECT id, title, created, edited, author, tags, description, showComment, visible, pinToTop FROM 'blog-post';")
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
			_tags       string
			description string
			showComment int
			visible     int
			pinToTop    int
		)
		if err := rows.Scan(&id, &title, &created, &edited, &author, &_tags, &description, &showComment, &visible, &pinToTop); err != nil {
			return nil, fmt.Errorf("Extracting data from result: %w", err)
		}
		var tags []string
		err := json.Unmarshal([]byte(_tags), &tags)
		if err != nil {
			return nil, fmt.Errorf("Unmarshal tags: %w", err)
		}
		post := PostMetadata{
			Id:          id,
			Title:       title,
			Created:     created,
			Edited:      edited,
			Author:      author,
			Tags:        tags,
			Description: description,
			ShowComment: showComment != 0,
			Visible:     visible != 0,
			PinToTop:    pinToTop != 0,
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (b Blog) GetPost(id string) (*PostMetadata, error) {
	row := b.db.QueryRow("SELECT title, created, edited, author, tags, description, showComment, visible, pinToTop FROM 'blog-post' WHERE id = ?", id)
	var (
		title       string
		created     string
		edited      string
		author      string
		_tags       string
		description string
		showComment int
		visible     int
		pinToTop    int
	)
	if err := row.Scan(&title, &created, &edited, &author, &_tags, &description, &showComment, &visible, &pinToTop); err != nil {
		return nil, fmt.Errorf("Extracting data from result: %w", err)
	}
	var tags []string
	err := json.Unmarshal([]byte(_tags), &tags)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal tags: %w", err)
	}
	post := PostMetadata{
		Id:          id,
		Title:       title,
		Created:     created,
		Edited:      edited,
		Author:      author,
		Tags:        tags,
		Description: description,
		ShowComment: showComment != 0,
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
				log.Printf("%s is reading blog post %s content.\n", r.RemoteAddr, id)
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
				log.Printf("%s is reading blog post %s.\n", r.RemoteAddr, id)
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
			log.Printf("%s is reading all blog posts.\n", r.RemoteAddr)
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
