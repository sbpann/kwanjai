package models

import (
	"kwanjai/config"
	"kwanjai/libraries"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// Board model.
type Board struct {
	UUID        string `json:"uuid"`
	User        string `json:"user"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type boardPerform interface {
	createBoard() (int, string, *Board)
	findBoard() (int, string, *Board)
	updateBoard() (int, string, *Board)
	deleteBoard() (int, string, *Board)
}

// NewBoard board method for interface with controller.
func NewBoard(perform boardPerform) (int, string, *Board) {
	status, message, board := perform.createBoard()
	return status, message, board
}

// FindBoard board method for interface with controller.
func FindBoard(perform boardPerform) (int, string, *Board) {
	status, message, board := perform.findBoard()
	return status, message, board
}

// UpdateBoard board method for interface with controller.
func UpdateBoard(perform boardPerform) (int, string, *Board) {
	status, message, board := perform.updateBoard()
	return status, message, board
}

// DeleteBoard board method for interface with controller.
func DeleteBoard(perform boardPerform) (int, string, *Board) {
	status, message, board := perform.deleteBoard()
	return status, message, board
}

func (board *Board) createBoard() (int, string, *Board) {
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	board.UUID = uuid.New().String()
	_, err = firestoreClient.Collection("boards").Doc(board.UUID).Set(config.Context, board)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusCreated, "Created board.", board
}

func (board *Board) findBoard() (int, string, *Board) {
	if board.UUID == "" {
		return http.StatusNotFound, "Board not found.", nil
	}
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	getBoard, err := firestoreClient.Collection("boards").Doc(board.UUID).Get(config.Context)
	if getBoard.Exists() {
		getBoard.DataTo(&board)
		return http.StatusOK, "Get board successfully.", board
	}
	return http.StatusNotFound, "Board not found.", nil
}

func (board *Board) updateBoard() (int, string, *Board) {
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	_, err = firestoreClient.Collection("boards").Doc(board.UUID).Update(config.Context, []firestore.Update{
		{
			Path:  "Name",
			Value: board.Name,
		},
		{
			Path:  "Description",
			Value: board.Description,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Updated sucessfully.", board
}

func (board *Board) deleteBoard() (int, string, *Board) {
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	_, err = firestoreClient.Collection("boards").Doc(board.UUID).Delete(config.Context)
	if err != nil {
		return http.StatusInternalServerError, err.Error(), nil
	}
	return http.StatusOK, "Deleted sucessfully.", nil
}
