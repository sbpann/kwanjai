package models

import (
	"kwanjai/libraries"
	"net/http"

	"github.com/google/uuid"
)

// Project model
type Project struct {
	UUID    string   `json:"uuid"`
	User    string   `json:"username"`
	Name    string   `json:"name" binding:"required"`
	Members []string `json:"members"`
}

func (project *Project) CreateProject() (int, string, *Project) {
	project.UUID = uuid.New().String()
	project.Members = append(project.Members, project.User)
	_, err := libraries.FirestoreCreatedOrSet("projects", project.UUID, project)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusCreated, "Created project.", project
}

func (project *Project) FindProject() (int, string, *Project) {
	if project.UUID == "" {
		return http.StatusNotFound, "Project not found.", nil
	}
	getProject, _ := libraries.FirestoreFind("projects", project.UUID)
	if getProject.Exists() {
		getProject.DataTo(project)
		return http.StatusOK, "Get project successfully.", project
	}
	return http.StatusNotFound, "Project not found.", nil
}

func (project *Project) UpdateProject() (int, string, *Project) {
	_, err := libraries.FirestoreUpdateField("projects", project.UUID, "Name", project.Name)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", project
}

func (project *Project) DeleteProject() (int, string, *Project) {
	_, err := libraries.FirestoreDelete("projects", project.UUID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
