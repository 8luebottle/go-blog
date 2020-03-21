package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/8luebottle/go-blog/api/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/8luebottle/go-blog/api/models"
	"github.com/8luebottle/go-blog/api/responses"
	"github.com/8luebottle/go-blog/api/utils/formaterror"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = json.Unmarshal(body, &user); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	if err = user.Validate("login"); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := s.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, token)
}

func (s *Server) SignIn(email, password string) (string, error) {
	user := models.User{}

	if err := s.DB.Model(models.User{}).
		Where("email = ?", email).
		Take(&user).Error; err != nil {
		return "", err
	}

	err := models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	return auth.CreateToken(user.ID)
}
