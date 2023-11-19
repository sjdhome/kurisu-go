package blog

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type BlogHandler struct {
	blog *Blog
}

func (h BlogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
	const REGEXP_ERROR = "Parsing regular expression for URL"
	onError := func(reason string, err error) {
		log.Println(fmt.Errorf("%s: %w.", reason, err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	if r.Method == "GET" {
		// HTTP GET
		allPosts, err := regexp.MatchString(`^\/blog\/post\/$`, r.URL.Path)
		if err != nil {
			onError(REGEXP_ERROR, err)
			return
		}
		if !allPosts {
			// Single post
			readContent, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z0-9\-]+\/content\/$`, r.URL.Path)
			if err != nil {
				onError(REGEXP_ERROR, err)
				return
			}
			readMetadata, err := regexp.MatchString(`^\/blog\/post\/[A-Za-z0-9\-]+\/$`, r.URL.Path)
			if err != nil {
				onError(REGEXP_ERROR, err)
				return
			}
			if readContent {
				// Read single post content
				re := regexp.MustCompile(`^\/blog\/post\/([A-Za-z0-9\-]+)\/content\/$`)
				id := string(re.FindSubmatch([]byte(r.URL.Path))[1])
				postContent, err := h.blog.getPostContent(id)
				if err != nil {
					onError(fmt.Sprintf("Reading blog post %s content", id), err)
					return
				}
				w.Header().Set("Content-Type", "text/markdown")
				w.Write([]byte(postContent))
			} else if readMetadata {
				// Get single post
				re := regexp.MustCompile(`^\/blog\/post\/([A-Za-z0-9\-]+)\/$`)
				id := string(re.FindSubmatch([]byte(r.URL.Path))[1])
				post, err := h.blog.getPost(id)
				if err != nil {
					onError(fmt.Sprintf("Reading blog post '%s'", id), err)
					return
				}
				j, err := json.Marshal(post)
				if err != nil {
					onError(fmt.Sprintf("Marshaling blog post '%s'", id), err)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(j)
			} else {
				// Invalid URL
				log.Printf("Invalid URL %s.\n", r.URL.Path)
				return
			}
		} else {
			// Get all posts
			posts, err := h.blog.getAllPosts()
			if err != nil {
				onError("Unable to get all blog posts", err)
				return
			}
			j, err := json.Marshal(posts)
			if err != nil {
				onError("Unable to marshal blog posts", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		}
	} else {
		// Not HTTP GET
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
