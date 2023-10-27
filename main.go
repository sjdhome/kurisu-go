package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"regexp"
)

type BlogHandler struct {
	blog *Blog
}

func (h BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		isSinglePost, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z\-]+\/`, r.URL.Path)
		if err != nil {
			log.Println("Unable to parse regular expression for URL.")
			log.Println(err)
			return
		}
		if isSinglePost {
			readContent, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z\-]+\/content\/$`, r.URL.Path)
			if err != nil {
				log.Println("Unable to parse regular expression for URL.")
				log.Println(err)
				return
			}
			readMetadata, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z\-]+\/$`, r.URL.Path)
			if err != nil {
				log.Println("Unable to parse regular expression for URL.")
				log.Println(err)
				return
			}
			if readContent {
				// Read single post content
				re := regexp.MustCompile(`\/blog\/post\/([A-Za-z\-]+)\/content\/$`)
				id := string(re.FindSubmatch([]byte(r.URL.Path))[1])
				postContent, err := h.blog.GetPostContent(id)
				log.Printf("%s is reading blog post %s content.\n", r.RemoteAddr, id)
				if err != nil {
					log.Printf("Unable to read blog post %s content.\n", id)
					log.Println(err)
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
					log.Printf("Unable to get blog post %s.\n", id)
					log.Println(err)
					return
				}
				j, err := json.Marshal(post)
				if err != nil {
					log.Printf("Unable to marshal blog post %s.\n", id)
					log.Println(err)
					return
				}
				w.Write(j)
			}
		} else if r.URL.Path == "/blog/post/" {
			// Get all posts
			log.Printf("%s is reading all blog posts.\n", r.RemoteAddr)
			posts, err := h.blog.GetAllPosts()
			if err != nil {
				log.Println("Unable to get all blog posts.")
				log.Println(err)
				return
			}
			j, err := json.Marshal(posts)
			if err != nil {
				log.Println("Unable to marshal blog posts.")
				log.Println(err)
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
		return nil, err
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
			return nil, err
		}
		var tags []string
		err := json.Unmarshal([]byte(_tags), &tags)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	var tags []string
	err := json.Unmarshal([]byte(_tags), &tags)
	if err != nil {
		return nil, err
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

var logFile string
var port int
var databaseFile string

func main() {
	flag.StringVar(&logFile, "log-file", "", "Specify log storage directory.")
	flag.IntVar(&port, "port", 3000, "HTTP server port.")
	flag.StringVar(&databaseFile, "database-file", "kurisu.sqlite3", "SQLite database file.")
	flag.Parse()
	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Unable to open log file.")
			log.Fatalln(err)
		}
		log.SetOutput(f)
	}

	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Println("Unable to open database.")
		log.Fatalln(err)
	}
	defer db.Close()

	blog := Blog{
		db: db,
	}
	blogHandler := BlogHandler{
		blog: &blog,
	}
	http.Handle("/blog/", blogHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	log.Printf("HTTP server started at :%d.\n", port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
