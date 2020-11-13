package models

import (
	"kwanjai/libraries"
	"log"
	"net/http"
	"time"
)

// Post model.
type Post struct {
	ID           string     `json:"id"`
	Board        string     `json:"board" binding:"required"`
	Project      string     `json:"project"` // It's an advantage for checking project membership.
	User         string     `json:"user"`
	Title        string     `json:"title" binding:"required"`
	Content      string     `json:"content" binding:"required"`
	Completed    bool       `json:"completed"`
	Urgent       bool       `json:"urgent"`
	Comments     []*Comment `json:"comments"`
	People       []string   `json:"people"`
	AddedDate    time.Time  `json:"added_date"`
	LastModified time.Time  `json:"last_modified"`
	DueDate      time.Time  `json:"due_date" binding:"required"`
}

// Comment model.
type Comment struct {
	UUID         string    `json:"uuid"`
	User         string    `json:"user"`
	Body         string    `json:"body"`
	AddedDate    time.Time `json:"added_date"`
	LastModified time.Time `json:"last_modified"`
}

func (post *Post) initialize() {
	post.People = []string{}
	post.Comments = []*Comment{}
	now := time.Now().Truncate(time.Millisecond)
	post.AddedDate = now
	post.LastModified = now
}

// CreatePost method returns status (int), message (string), post object.
func (post *Post) CreatePost() (int, string, *Post) {
	post.initialize()
	reference, _, err := libraries.FirestoreAdd("posts", post)
	if err != nil {
		log.Panicln(err)
	}
	go libraries.FirestoreDeleteField("posts", reference.ID, "ID")
	post.ID = reference.ID
	return http.StatusCreated, "Created post.", post
}

// FindPost method returns status (int), message (string), post object.
func (post *Post) FindPost() (int, string, *Post) {
	if post.ID == "" {
		return http.StatusNotFound, "Post not found.", nil
	}
	getPost, _ := libraries.FirestoreFind("posts", post.ID)
	if getPost.Exists() {
		getPost.DataTo(post)
		post.ID = getPost.Ref.ID
		return http.StatusOK, "Get post successfully.", post
	}
	return http.StatusNotFound, "Post not found.", nil
}

// UpdatePost method returns status (int), message (string), post object.
func (post *Post) UpdatePost(field string, value interface{}) (int, string, *Post) {
	_, err := libraries.FirestoreUpdateField("posts", post.ID, field, value)
	if err != nil {
		log.Panicln(err)
	}
	return http.StatusOK, "Updated sucessfully.", post
}

// DeletePost method returns status (int), message (string), nil object.
func (post *Post) DeletePost() (int, string, *Post) {
	_, err := libraries.FirestoreDelete("posts", post.ID)
	if err != nil {
		log.Panicln(err)
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
