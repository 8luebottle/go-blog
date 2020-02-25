package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/8luebottle/go-blog/api/models"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreatePost(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password")
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		authorId     uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"ConfigMaps and Secrets on a Gorush server example", "content:": "We have a Gorush server from the Kubernetes: running a push-server with Gorush behind an AWS LoadBalancer post, and I’d like to add an ability to configure it via a Github repository and run it with a different set of settings — for Staging and Production environments.", "authorId": 1`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "ConfigMaps and Secrets on a Gorush server example",
			content:      "We have a Gorush server from the Kubernetes: running a push-server with Gorush behind an AWS LoadBalancer post, and I’d like to add an ability to configure it via a Github repository and run it with a different set of settings — for Staging and Production environments.",
			authorId:     user.ID,
			errorMessage: "",
		},
		{
			inputJSON:    `{"title": "Feb 18, 2020 05:17:20AM", "content:":"Early Bird", "authorId": 2}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "ConfigMaps and Secrets on a Gorush server example", "content:":"Blah Blah.", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			inputJSON:    `{"title": "Daily Commit", "content":"Consistency", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "Incorrect Token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content":"Empty Title", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Title Require",
		},
		{
			inputJSON:    `{"title":"The Title", "content": "The Content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Author Require",
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreatePost)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["authorId"], float64(v.authorId))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetPosts(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPosts)
	handler.ServeHTTP(rr, req)

	var posts []models.Post
	err = json.Unmarshal([]byte(rr.Body.String()), &posts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(posts), 2)
}

func TestGetPostByID(t *testing.T) {
	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	post, err := seedOneUserAndOnePost()
	if err != nil {
		log.Fatal(err)
	}
	postSample := []struct {
		id           string
		statusCode   int
		title        string
		content      string
		authorId     uint32
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(post.ID)),
			statusCode: 200,
			title:      post.Title,
			content:    post.Content,
			authorId:   post.AuthorID,
		},
		{
			id:         "anonymous",
			statusCode: 400,
		},
	}
	for _, v := range postSample {
		req, err := http.NewRequest("GET", "/posts", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetPost)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, post.Title, responseMap["title"])
			assert.Equal(t, post.Content, responseMap["content"])
			assert.Equal(t, float64(post.AuthorID), responseMap["authorId"])
		}
	}
}

func TestUpdatePost(t *testing.T) {
	var PostUserEmail, PostUserPassword string
	var AuthPostAuthorID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	// Get the first user
	for _, user := range users {
		if user.ID == 2 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password"
	}
	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		AuthPostID = post.ID
		AuthPostAuthorID = post.AuthorID
	}

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		content      string
		authorId     uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title": "CORONA19", "contet": "More than 80 clinical trials launch to test coronavirus treatments", "authorId":2}`,
			statusCode:   200,
			title:        "CORONA19",
			content:      "More than 80 clinical trials launch to test coronavirus treatments",
			authorId:     AuthPostAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			updateJSON:   `{"title":"UPDATED TITLE", "content": "UPDATED CONTENT", "authorId":1}`,
			tokenGiven:   "Wrong Token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			statusCode: 400,
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdatePost)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["authorId"], float64(v.authorId))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeletePost(t *testing.T) {
	var PostUserEmail, PostUserPassword string
	var PostUserID uint32
	var AuthPostID uint64

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}
	users, posts, err := seedUsersAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	// Get the second user
	for _, user := range users {
		if user.ID == 1 {
			continue
		}
		PostUserEmail = user.Email
		PostUserPassword = "password"
	}

	token, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	for _, post := range posts {
		if post.ID == 1 {
			continue
		}
		AuthPostID = post.ID
		PostUserID = post.AuthorID
	}
	postSample := []struct {
		id           string
		authorId     uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int 64 --> int before convert to string
			id:           strconv.Itoa(int(AuthPostID)),
			authorId:     PostUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// Empty Token
			id:           strconv.Itoa(int(AuthPostID)),
			authorId:     PostUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(AuthPostID)),
			authorId:     PostUserID,
			tokenGiven:   "Passed Incorrect Token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknown",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(1)),
			authorId:     1,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range postSample {
		req, _ := http.NewRequest("GET", "/posts", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeletePost)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
