package models

import (
	"kwanjai/libraries"
	"net/http"

	"github.com/google/uuid"
)

// Board model.
type Board struct {
	UUID        string `json:"uuid"`
	User        string `json:"user"`
	Name        string `json:"name" binding:"required"`
	Project     string `json:"project" binding:"required"`
	Description string `json:"description"`
}

func (board *Board) CreateBoard() (int, string, *Board) {
	board.UUID = uuid.New().String()
	_, err := libraries.FirestoreCreatedOrSet("boards", board.UUID, board)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusCreated, "Created board.", board
}

func (board *Board) FindBoard() (int, string, *Board) {
	if board.UUID == "" {
		return http.StatusNotFound, "Board not found.", nil
	}
	getBoard, _ := libraries.FirestoreFind("boards", board.UUID)
	if getBoard.Exists() {
		getBoard.DataTo(board)
		return http.StatusOK, "Get board successfully.", board
	}
	return http.StatusNotFound, "Board not found.", nil
}

// UpdateBoard board method.
func (board *Board) UpdateBoard(field string, value interface{}) (int, string, *Board) {
	_, err := libraries.FirestoreUpdateField("boards", board.UUID, field, value)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", board
}

func (board *Board) DeleteBoard() (int, string, *Board) {
	_, err := libraries.FirestoreDelete("boards", board.UUID)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
