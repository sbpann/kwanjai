package models

// Post model.
type Post struct {
	UUID         string   `json:"uuid"`
	Board        string   `json:"board"`
	User         string   `json:"username"`
	Title        string   `json:"title"`
	Body         string   `json:"body"`
	Completed    bool     `json:"is_completed"`
	Urgent       bool     `json:"is_urgent"`
	People       []string `json:"people"`
	AddedDate    string   `json:"added_date"`
	LastModified string   `json:"last_modified"`
}

// Comment model.
type Comment struct {
	UUID         string `json:"uuid"`
	User         string `json:"username"`
	Post         string `json:"title"`
	Body         string `json:"body"`
	AddedDate    string `json:"added_date"`
	LastModified string `json:"last_modified"`
}

// Reply model.
type Reply struct {
	UUID         string `json:"uuid"`
	User         string `json:"username"`
	Post         string `json:"title"`
	Body         string `json:"body"`
	AddedDate    string `json:"added_date"`
	LastModified string `json:"last_modified"`
}
