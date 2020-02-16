package controllertests

import (
	"errors"
	"fmt"
	"log"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("this is the error %v\n", err)
	}
	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     user.Nickname + "123",
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {
		token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
		}
	}
}

//
//func TestLogin(t *testing.T){
//	refreshUserTable()
//	_, err := seedOneUser()
//	if err != nil {
//		fmt.Printf("this is the error %v\n", err)
//	}
//	samples := []struct {
//		inputJSON string
//		statusCode int
//		email string
//		password string
//		errorMessage string
//	}{
//		inputJSON:i // Keep Writing from here
//	}
//}
