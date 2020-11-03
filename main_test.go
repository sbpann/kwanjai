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

func clearTestUser1(t *testing.T) {

	getUser, err := libraries.FirestoreFind("users", "test1")
	if getUser.Exists() {
		_, err = libraries.FirestoreDelete("users", "test1")
		assert.Equal(t, nil, err)
		getToken, err := libraries.FirestoreSearch("tokenUUID", "user", "==", "test1")
		assert.Equal(t, nil, err)
		for _, token := range getToken {
			_, err = libraries.FirestoreDelete("tokenUUID", token.Ref.ID)
			assert.Equal(t, nil, err)
		}
	}
	getEmail, err := libraries.FirestoreSearch("users", "Email", "==", "test1@example.com")
	assert.Equal(t, nil, err)
	if len(getEmail) > 0 {
		_, err = libraries.FirestoreDelete("users", getEmail[0].Data()["Username"].(string))
		assert.Equal(t, nil, err)
		_, err = libraries.FirestoreDelete("tokenUUID", getEmail[0].Data()["Username"].(string))
		assert.Equal(t, nil, err)
		getToken, err := libraries.FirestoreSearch("tokenUUID", "user", "==", getEmail[0].Data()["Username"])
		assert.Equal(t, nil, err)
		for _, token := range getToken {
			_, err = libraries.FirestoreDelete("tokenUUID", token.Ref.ID)
			assert.Equal(t, nil, err)
		}
	}
}

func clearTestUser2(t *testing.T) {
	getUser, err := libraries.FirestoreFind("users", "test2")
	if getUser.Exists() {
		_, err = libraries.FirestoreDelete("users", "test2")
		assert.Equal(t, nil, err)
		getToken, err := libraries.FirestoreSearch("tokenUUID", "user", "==", "test2")
		assert.Equal(t, nil, err)
		for _, token := range getToken {
			_, err = libraries.FirestoreDelete("tokenUUID", token.Ref.ID)
			assert.Equal(t, nil, err)
		}
	}
	getEmail, err := libraries.FirestoreSearch("users", "Email", "==", "test2@example.com")
	assert.Equal(t, nil, err)
	if len(getEmail) > 0 {
		_, err = libraries.FirestoreDelete("users", getEmail[0].Data()["Username"].(string))
		assert.Equal(t, nil, err)
		_, err = libraries.FirestoreDelete("tokenUUID", getEmail[0].Data()["Username"].(string))
		assert.Equal(t, nil, err)
		getToken, err := libraries.FirestoreSearch("tokenUUID", "user", "==", getEmail[0].Data()["Username"])
		assert.Equal(t, nil, err)
		for _, token := range getToken {
			_, err = libraries.FirestoreDelete("tokenUUID", token.Ref.ID)
			assert.Equal(t, nil, err)
		}
	}
}

func TestSetup(t *testing.T) {
	setupServer()
	clearTestUser1(t)
	clearTestUser2(t)
}

func TestRigesterLogoutLoginLogout(t *testing.T) {
	clearTestUser1(t)
	// register
	registerInfo := new(models.User)
	registerInfo.Username = "test1"
	registerInfo.Email = "test1@example.com"
	registerInfo.Password = "test1password"
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
	login.ID = "test1"
	login.Password = "test1password"
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

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/verify_email/badlink", nil)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusBadRequest, writer.Code)
	assert.Equal(t, `{"message":"Bad verification link."}`, writer.Body.String())
}

func TestUnauthorizedAction(t *testing.T) {

	endpoints := map[string]string{
		"/project/new":    "POST",
		"/project/find":   "POST",
		"/project/update": "PATCH",
		"/project/delete": "DELETE",
	}
	for key, element := range endpoints {
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(element, key, nil)
		getServer("test").ServeHTTP(writer, request)
		assert.Equal(t, http.StatusUnauthorized, writer.Code)
	}
}
func TestCreateBoardGetAllProjects(t *testing.T) {
	clearTestUser1(t)
	clearTestUser2(t)
	var response map[string]interface{}

	// register1
	registerInfo := new(models.User)
	registerInfo.Username = "test1"
	registerInfo.Email = "test1@example.com"
	registerInfo.Password = "test1password"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token1 := new(libraries.Token)
	token1.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// register2
	registerInfo = new(models.User)
	registerInfo.Username = "test2"
	registerInfo.Email = "test2@example.com"
	registerInfo.Password = "test2password"
	b, _ = json.Marshal(registerInfo)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token2 := new(libraries.Token)
	token2.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// created project by user test1
	project := new(models.Project)
	project.Name = "My New Project"
	b, _ = json.Marshal(project)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/project/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdProjectUUID := response["project"].(map[string]interface{})["uuid"].(string)

	// user test2 get all projects
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/project/all", nil)
	request.Header.Set("Authorization", token2.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	allProjects := response["projects"].([]interface{})
	assert.Equal(t, 0, len(allProjects)) // should return blank array

	// user test1 get all projects
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/project/all", nil)
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	allProjects = response["projects"].([]interface{})
	assert.Equal(t, 1, len(allProjects)) //should return array with one element

	_, err := libraries.FirestoreDelete("projects", createdProjectUUID)
	assert.Equal(t, nil, err)
}

func TestCreateBoard(t *testing.T) {
	clearTestUser1(t)
	clearTestUser2(t)
	var response map[string]interface{}

	// register1
	registerInfo := new(models.User)
	registerInfo.Username = "test1"
	registerInfo.Email = "test1@example.com"
	registerInfo.Password = "test1password"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token1 := new(libraries.Token)
	token1.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// register2
	registerInfo = new(models.User)
	registerInfo.Username = "test2"
	registerInfo.Email = "test2@example.com"
	registerInfo.Password = "test2password"
	b, _ = json.Marshal(registerInfo)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token2 := new(libraries.Token)
	token2.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// created project by user test1
	project := new(models.Project)
	project.Name = "My New Project"
	b, _ = json.Marshal(project)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/project/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdProjectUUID := response["project"].(map[string]interface{})["uuid"].(string)

	// user test2 try to create board under project of user test1.
	board := new(models.Board)
	board.Name = "My new board"
	board.Project = createdProjectUUID
	b, _ = json.Marshal(board)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/board/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token2.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusForbidden, writer.Code) // should fail.

	// user test1 try to create board under project of user test1.
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/board/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code) // should pass.
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdBoardtUUID := response["board"].(map[string]interface{})["uuid"].(string)

	_, err := libraries.FirestoreDelete("projects", createdProjectUUID)
	assert.Equal(t, nil, err)
	_, err = libraries.FirestoreDelete("boards", createdBoardtUUID)
	assert.Equal(t, nil, err)
}

func TestCreateAndDeletePost(t *testing.T) {
	clearTestUser1(t)
	clearTestUser2(t)
	var response map[string]interface{}

	// register1
	registerInfo := new(models.User)
	registerInfo.Username = "test1"
	registerInfo.Email = "test1@example.com"
	registerInfo.Password = "test1password"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token1 := new(libraries.Token)
	token1.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// Created project
	project := new(models.Project)
	project.Name = "My New Project"
	b, _ = json.Marshal(project)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/project/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdProjectUUID := response["project"].(map[string]interface{})["uuid"].(string)

	// Created board
	board := new(models.Board)
	board.Name = "My new board"
	board.Project = createdProjectUUID
	b, _ = json.Marshal(board)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/board/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code) // should pass.
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdBoardtUUID := response["board"].(map[string]interface{})["uuid"].(string)

	// Created post
	post := new(models.Post)
	post.Board = createdBoardtUUID
	post.Title = "My post"
	post.Body = "My this post is created for testing."
	b, _ = json.Marshal(post)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/post/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code) // should pass.
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdPostUUID := response["post"].(map[string]interface{})["uuid"].(string)

	// register2
	registerInfo = new(models.User)
	registerInfo.Username = "test2"
	registerInfo.Email = "test2@example.com"
	registerInfo.Password = "test2password"
	b, _ = json.Marshal(registerInfo)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token2 := new(libraries.Token)
	token2.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// user test2 try to delete post created by user 1.
	post.UUID = createdPostUUID
	b, _ = json.Marshal(post)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/post/delete", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token2.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusForbidden, writer.Code) // should fail.

	// user test1 try to delete post created by user 1.
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/post/delete", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code) // should pass.

	_, err := libraries.FirestoreDelete("projects", createdProjectUUID)
	assert.Equal(t, nil, err)
	_, err = libraries.FirestoreDelete("boards", createdBoardtUUID)
	assert.Equal(t, nil, err)
}

func TestCreateAndDeleteComment(t *testing.T) {
	clearTestUser1(t)
	clearTestUser2(t)
	var response map[string]interface{}

	// register1
	registerInfo := new(models.User)
	registerInfo.Username = "test1"
	registerInfo.Email = "test1@example.com"
	registerInfo.Password = "test1password"
	b, _ := json.Marshal(registerInfo)
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token1 := new(libraries.Token)
	token1.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// Created project
	project := new(models.Project)
	project.Name = "My New Project"
	b, _ = json.Marshal(project)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/project/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdProjectUUID := response["project"].(map[string]interface{})["uuid"].(string)

	// Created board
	board := new(models.Board)
	board.Name = "My new board"
	board.Project = createdProjectUUID
	b, _ = json.Marshal(board)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/board/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code) // should pass.
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdBoardtUUID := response["board"].(map[string]interface{})["uuid"].(string)

	// Created post
	post := new(models.Post)
	post.Board = createdBoardtUUID
	post.Title = "My post"
	post.Body = "My this post is created for testing."
	b, _ = json.Marshal(post)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/post/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code) // should pass.
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdPostUUID := response["post"].(map[string]interface{})["uuid"].(string)
	post.UUID = createdPostUUID

	// Create comment by user test1
	post.Comments = append(post.Comments, &models.Comment{Body: "my comment"})
	b, _ = json.Marshal(post)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/post/comment/new", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusCreated, writer.Code) // should pass.
	json.Unmarshal([]byte(writer.Body.String()), &response)
	createdCommentUUID := response["post"].(map[string]interface{})["comments"].([]interface{})[0].(map[string]interface{})["uuid"].(string)
	post.Comments[0].UUID = createdCommentUUID

	// register2
	registerInfo = new(models.User)
	registerInfo.Username = "test2"
	registerInfo.Email = "test2@example.com"
	registerInfo.Password = "test2password"
	b, _ = json.Marshal(registerInfo)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(b)))
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code)
	json.Unmarshal([]byte(writer.Body.String()), &response)
	token2 := new(libraries.Token)
	token2.AccessToken = response["token"].(map[string]interface{})["access_token"].(string)

	// user2 try to delete commented by user1
	b, _ = json.Marshal(post)
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/post/comment/delete", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token2.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusForbidden, writer.Code) // should fail.

	// user1 try to delete commented by user1
	writer = httptest.NewRecorder()
	request, _ = http.NewRequest("DELETE", "/post/comment/delete", bytes.NewBuffer([]byte(b)))
	request.Header.Set("Authorization", token1.AccessToken)
	getServer("test").ServeHTTP(writer, request)
	assert.Equal(t, http.StatusOK, writer.Code) // should pass.

	_, err := libraries.FirestoreDelete("projects", createdProjectUUID)
	assert.Equal(t, nil, err)
	_, err = libraries.FirestoreDelete("boards", createdBoardtUUID)
	assert.Equal(t, nil, err)
	_, err = libraries.FirestoreDelete("posts", createdPostUUID)
	assert.Equal(t, nil, err)
}
