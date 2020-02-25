package controllers

import (
	"net/http"

	"github.com/8luebottle/go-blog/api/responses"
)

// Welcome API User
func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "welcome to this Awesome API")
}
