package models

import (
	"kwanjai/libraries"
	"net/http"
)

// Board model.
type Board struct {
	ID          string `json:"id"`
	User        string `json:"user"`
	Name        string `json:"name" binding:"required"`
	Project     string `json:"project" binding:"required"`
	Description string `json:"description"`
}

func (board *Board) CreateBoard() (int, string, *Board) {
	reference, _, err := libraries.FirestoreAdd("boards", board)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	go libraries.FirestoreDeleteField("boards", reference.ID, "ID")
	board.ID = reference.ID
	return http.StatusCreated, "Created board.", board
}

func (board *Board) FindBoard() (int, string, *Board) {
	if board.ID == "" {
		return http.StatusNotFound, "Board not found.", nil
	}
	getBoard, _ := libraries.FirestoreFind("boards", board.ID)
	if getBoard.Exists() {
		getBoard.DataTo(board)
		board.ID = getBoard.Ref.ID
		return http.StatusOK, "Get board successfully.", board
	}
	return http.StatusNotFound, "Board not found.", nil
}

// UpdateBoard board method.
func (board *Board) UpdateBoard(field string, value interface{}) (int, string, *Board) {
	_, err := libraries.FirestoreUpdateField("boards", board.ID, field, value)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", board
}

func (board *Board) DeleteBoard() (int, string, *Board) {
	_, err := libraries.FirestoreDelete("boards", board.ID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
