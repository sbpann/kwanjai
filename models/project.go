package models

import (
	"kwanjai/libraries"
	"log"
	"net/http"
	"time"
)

// Project model
type Project struct {
	ID          string    `json:"id"`
	User        string    `json:"user"`
	Name        string    `json:"name" binding:"required"`
	Members     []string  `json:"members"`
	CreatedDate time.Time `json:"created_date"`
}

// CreateProject method returns status (int), message (string), project object.
func (project *Project) CreateProject() (int, string, *Project) {
	project.Members = append(project.Members, project.User)
	project.CreatedDate = time.Now().Truncate(time.Microsecond)
	reference, _, err := libraries.FirestoreAdd("projects", project)
	if err != nil {
		log.Panicln(err)
	}
	go libraries.FirestoreDeleteField("projects", reference.ID, "ID")
	project.ID = reference.ID
	return http.StatusCreated, "Created project.", project
}

// FindProject method returns status (int), message (string), project object.
func (project *Project) FindProject() (int, string, *Project) {
	if project.ID == "" {
		return http.StatusNotFound, "Project not found.", nil
	}
	getProject, _ := libraries.FirestoreFind("projects", project.ID)
	if getProject.Exists() {
		getProject.DataTo(project)
		project.ID = getProject.Ref.ID
		return http.StatusOK, "Get project successfully.", project
	}
	return http.StatusNotFound, "Project not found.", nil
}

// UpdateProject method returns status (int), message (string), project object.
func (project *Project) UpdateProject() (int, string, *Project) {
	_, err := libraries.FirestoreUpdateField("projects", project.ID, "Name", project.Name)
	if err != nil {
		log.Panicln(err)
	}
	_, err = libraries.FirestoreUpdateField("projects", project.ID, "Members", project.Members)
	if err != nil {
		log.Panicln(err)
	}
	return http.StatusOK, "Updated sucessfully.", project
}

// DeleteProject method returns status (int), message (string), nil object.
func (project *Project) DeleteProject() (int, string, *Project) {
	_, err := libraries.FirestoreDelete("projects", project.ID)
	if err != nil {
		log.Panicln(err)
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
