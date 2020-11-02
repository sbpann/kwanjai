package models

import (
	"kwanjai/libraries"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Post model.
type Post struct {
	UUID         string     `json:"uuid"`
	Board        string     `json:"board" binding:"required"`
	User         string     `json:"username"`
	Title        string     `json:"title" binding:"required"`
	Body         string     `json:"body" binding:"required"`
	Comments     []*Comment `json:"comments"`
	Completed    bool       `json:"is_completed"`
	Urgent       bool       `json:"is_urgent"`
	People       []string   `json:"people"`
	AddedDate    time.Time  `json:"added_date"`
	LastModified time.Time  `json:"last_modified"`
}

// Comment model.
type Comment struct {
	UUID         string    `json:"uuid"`
	User         string    `json:"username"`
	Body         string    `json:"body"`
	AddedDate    time.Time `json:"added_date"`
	LastModified time.Time `json:"last_modified"`
}

func (post *Post) initialize() {
	post.UUID = uuid.New().String()
	post.People = []string{}
	now := time.Now().Truncate(time.Millisecond)
	post.Comments = []*Comment{}
	post.AddedDate = now
	post.LastModified = now
}

func (post *Post) CreatePost() (int, string, *Post) {
	post.initialize()
	_, err := libraries.FirestoreCreatedOrSet("posts", post.UUID, post)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusCreated, "Created post.", post
}

func (post *Post) FindPost() (int, string, *Post) {
	if post.UUID == "" {
		return http.StatusNotFound, "Post not found.", nil
	}
	getPost, _ := libraries.FirestoreFind("posts", post.UUID)
	if getPost.Exists() {
		getPost.DataTo(post)
		return http.StatusOK, "Get post successfully.", post
	}
	return http.StatusNotFound, "Post not found.", nil
}

func (post *Post) UpdatePost(field string, value interface{}) (int, string, *Post) {
	_, err := libraries.FirestoreUpdateField("posts", post.UUID, field, value)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", post
}

func (post *Post) DeletePost() (int, string, *Post) {
	_, err := libraries.FirestoreDelete("posts", post.UUID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
