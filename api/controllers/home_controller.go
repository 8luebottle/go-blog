package controllers

import (
	"net/http"

	"github.com/8luebottle/go-blog/api/responses"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "GO-Blog API v.1")
}
