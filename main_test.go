package main

import (
	"bytes"
	"encoding/json"
	"kwanjai/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

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
	assert.Equal(t, 400, writer.Code)
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
	assert.Equal(t, 400, writer.Code)
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
	assert.Equal(t, 400, writer.Code)
	assert.Equal(t, `{"message":"Cannot login with provided credential."}`, writer.Body.String())
}
