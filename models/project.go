package models

import (
	"kwanjai/libraries"
	"net/http"
)

// Project model
type Project struct {
	ID      string   `json:"id"`
	User    string   `json:"username"`
	Name    string   `json:"name" binding:"required"`
	Members []string `json:"members"`
}

func (project *Project) CreateProject() (int, string, *Project) {
	project.Members = append(project.Members, project.User)
	reference, _, err := libraries.FirestoreAdd("projects", project)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	go libraries.FirestoreDeleteField("projects", reference.ID, "ID")
	project.ID = reference.ID
	return http.StatusCreated, "Created project.", project
}

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

func (project *Project) UpdateProject() (int, string, *Project) {
	_, err := libraries.FirestoreUpdateField("projects", project.ID, "Name", project.Name)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", project
}

func (project *Project) DeleteProject() (int, string, *Project) {
	_, err := libraries.FirestoreDelete("projects", project.ID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
