package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

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
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"ConfigMaps and Secrets on a Gorush server example", "content:": "We have a Gorush server from the Kubernetes: running a push-server with Gorush behind an AWS LoadBalancer post, and I’d like to add an ability to configure it via a Github repository and run it with a different set of settings — for Staging and Production environments.", "author_id": 1`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "ConfigMaps and Secrets on a Gorush server example",
			content:      "We have a Gorush server from the Kubernetes: running a push-server with Gorush behind an AWS LoadBalancer post, and I’d like to add an ability to configure it via a Github repository and run it with a different set of settings — for Staging and Production environments.",
			author_id:    user.ID,
			errorMessage: "",
		}, // Write further Testing here
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/posts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handeler := http.HandlerFunc(server.CreatePost)

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
			assert.Equal(t, responseMap["author_id"], float64(v.author_id))
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

// func TestGetPosts
