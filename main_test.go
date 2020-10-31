package main

import (
	"bytes"
	"encoding/json"
	"kwanjai/config"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/go-playground/assert/v2"
)

func clearTestUser(t *testing.T, firestoreClient *firestore.Client) {

	_, err := firestoreClient.Collection("users").Doc("test").Get(config.Context)
	if err == nil {
		_, err = firestoreClient.Collection("users").Doc("test").Delete(config.Context)
		assert.Equal(t, nil, err)
	}
	findEmail := firestoreClient.Collection("users").Where("Email", "==", "test@example.com").Documents(config.Context)
	foundEmail, err := findEmail.GetAll()
	assert.Equal(t, nil, err)
	if len(foundEmail) > 0 {
		user := new(models.User)
		_ = foundEmail[0].DataTo(&user)
		_, err = firestoreClient.Collection("users").Doc(user.Username).Delete(config.Context)
		assert.Equal(t, nil, err)
		_, err = firestoreClient.Collection("tokenUUID").Doc(user.Username).Delete(config.Context)
		assert.Equal(t, nil, err)
	}
}

func TestRegisterWithAGoodinfo(t *testing.T) {
	setupServer()
	// find test user and delete
	firestoreClient, err := libraries.FirebaseApp().Firestore(config.Context)
	defer firestoreClient.Close()
	assert.Equal(t, nil, err)
	clearTestUser(t, firestoreClient)
	// find test user and delete

	registerInfo := new(models.User)
	registerInfo.Username = "test"
	registerInfo.Email = "test@example.com"
	registerInfo.Password = "testpassword"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, 200, writer.Code)

	var result map[string]interface{}
	json.Unmarshal([]byte(writer.Body.String()), &result)
	assert.Equal(t, result["warning"].(string), "You have just registered with the username (test) or the email (test@example.com) which is going to be delete eventually. Please avoid using those names.")

	// clear data
	_, err = firestoreClient.Collection("users").Doc("test").Delete(config.Context)
	assert.Equal(t, nil, err)
	_, err = firestoreClient.Collection("tokenUUID").Doc("test").Delete(config.Context)
	assert.Equal(t, nil, err)
}

// func TestRegisterWithBadEmailFormat(t *testing.T) {
// 	setupServer()
// 	registerInfo := new(models.User)
// 	registerInfo.Username = "john"
// 	registerInfo.Email = "bad-email"
// 	registerInfo.Password = "johnpassword"
// 	b, _ := json.Marshal(registerInfo)
// 	writer := httptest.NewRecorder()
// 	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
// 	getServer("test").ServeHTTP(writer, request)
// 	assert.Equal(t, 400, writer.Code)
// }

// func TestRegisterWithReserverdUsername(t *testing.T) {
// 	setupServer()
// 	registerInfo := new(models.User)
// 	registerInfo.Username = "anonymous"
// 	registerInfo.Email = "anonymous@email.com"
// 	registerInfo.Password = "anonymouspassword"
// 	b, _ := json.Marshal(registerInfo)
// 	writer := httptest.NewRecorder()
// 	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
// 	getServer("test").ServeHTTP(writer, request)
// 	assert.Equal(t, 400, writer.Code)
// }

// func TestLoginWithInvalidCredential(t *testing.T) {
// 	setupServer()
// 	login := new(models.LoginCredential)
// 	login.ID = "anonymous"
// 	login.Password = "anonymouspassword"
// 	b, _ := json.Marshal(login)
// 	writer := httptest.NewRecorder()
// 	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(b)))
// 	getServer("test").ServeHTTP(writer, request)
// 	assert.Equal(t, 400, writer.Code)
// 	assert.Equal(t, `{"message":"Cannot login with provided credential."}`, writer.Body.String())
// }

// func TestUnauthorizedBoardAction(t *testing.T) {
// 	setupServer()
// 	board := new(models.Board)
// 	board.Name = "myboardname"
// 	b, _ := json.Marshal(board)
// 	endpoints := map[string]string{
// 		"/board/new":    "POST",
// 		"/board/find":   "GET",
// 		"/board/update": "PATCH",
// 		"/board/delete": "DELETE",
// 	}
// 	for key, element := range endpoints {
// 		writer := httptest.NewRecorder()
// 		request, _ := http.NewRequest(element, key, bytes.NewBuffer([]byte(b)))
// 		getServer("test").ServeHTTP(writer, request)
// 		assert.Equal(t, http.StatusUnauthorized, writer.Code)
// 		assert.Equal(t, `{"message":"Authenticated user only."}`, writer.Body.String())
// 	}
// }
