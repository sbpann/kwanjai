package main

import (
	"bytes"
	"encoding/json"
	"kwanjai/libraries"
	"kwanjai/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func clearTestUser(t *testing.T) {

	getUser, err := libraries.FirestoreFind("users", "test")
	if getUser.Exists() {
		_, err = libraries.FirestoreDelete("users", "test")
		assert.Equal(t, nil, err)
		_, err = libraries.FirestoreDelete("tokenUUID", "test")
		assert.Equal(t, nil, err)
	}
	getEmail, err := libraries.FirestoreSearch("users", "Email", "==", "test@example.com")
	assert.Equal(t, nil, err)
	if len(getEmail) > 0 {
		_, err = libraries.FirestoreDelete("users", getEmail[0].Data()["Username"].(string))
		assert.Equal(t, nil, err)
		_, err = libraries.FirestoreDelete("tokenUUID", getEmail[0].Data()["Username"].(string))
		assert.Equal(t, nil, err)
	}
}

func TestRegisterWithAGoodInfo(t *testing.T) {
	setupServer()
	clearTestUser(t)

	registerInfo := new(models.User)
	registerInfo.Username = "test"
	registerInfo.Email = "test@example.com"
	registerInfo.Password = "testpassword"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)

	var response map[string]interface{}
	json.Unmarshal([]byte(writer.Body.String()), &response)
	assert.Equal(t, response["warning"].(string), "You have just registered with the username (test) or the email (test@example.com) which is going to be delete eventually. Please avoid using those names.")
}

func TestRigesterLogoutLoginLogout(t *testing.T) {
	setupServer()
	clearTestUser(t)

	// register
	registerInfo := new(models.User)
	registerInfo.Username = "test"
	registerInfo.Email = "test@example.com"
	registerInfo.Password = "testpassword"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)

	// Logout
	var response map[string]interface{}
	json.Unmarshal([]byte(writer.Body.String()), &response)
	writer = httptest.NewRecorder()
	token := new(libraries.Token)
	token.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)
	token.RefreshToken = response["token"].(map[string]interface{})["refresh_token"].(string)
	b, _ = json.Marshal(token)
	request, _ = http.NewRequest("POST", "/logout", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	assert.Equal(t, response["message"].(string), "User logged out successfully.")

	//Login
	writer = httptest.NewRecorder()
	login := new(models.LoginCredential)
	login.ID = "test"
	login.Password = "testpassword"
	b, _ = json.Marshal(login)
	request, _ = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)

	// Logout
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token = new(libraries.Token)
	token.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)
	token.RefreshToken = response["token"].(map[string]interface{})["refresh_token"].(string)
	writer = httptest.NewRecorder()
	b, _ = json.Marshal(token)
	request, _ = http.NewRequest("POST", "/logout", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	assert.Equal(t, response["message"].(string), "User logged out successfully.")
}

func TestRegisterWithBadEmailFormat(t *testing.T) {
	setupServer()
	registerInfo := new(models.User)
	registerInfo.Username = "john"
	registerInfo.Email = "bad-email"
	registerInfo.Password = "johnpassword"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestRegisterWithReserverdUsername(t *testing.T) {
	setupServer()
	registerInfo := new(models.User)
	registerInfo.Username = "anonymous"
	registerInfo.Email = "anonymous@email.com"
	registerInfo.Password = "anonymouspassword"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestLoginWithInvalidCredential(t *testing.T) {
	setupServer()
	login := new(models.LoginCredential)
	login.ID = "anonymous"
	login.Password = "anonymouspassword"
	b, _ := json.Marshal(login)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, `{"message":"Cannot login with provided credential."}`, writer.Body.String())
}

func TestVerifyEmailWithBadLink(t *testing.T) {
	setupServer()
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/verify_email/badlink", nil)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, `{"message":"Bad verification link."}`, writer.Body.String())
}

func TestUnauthorizedBoardAction(t *testing.T) {
	setupServer()
	board := new(models.Board)
	board.Name = "myboardname"
	b, _ := json.Marshal(board)
	endpoints := map[string]string{
		"/board/new":    "POST",
		"/board/find":   "POST",
		"/board/update": "PATCH",
		"/board/delete": "DELETE",
	}
	for key, element := range endpoints {
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(element, key, bytes.NewBuffer([]byte(b)))
		getServer("test").ServeHTTP(writer, request)
		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	}
}
