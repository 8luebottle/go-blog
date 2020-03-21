package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/8luebottle/go-blog/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (s *Server) Initialize(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	if DbDriver == "mysql" {
		DBURL := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			DbUser, DbPassword, DbHost, DbPort, DbName)
		s.DB, err = gorm.Open(DbDriver, DBURL)
		if err != nil {
			fmt.Printf("cannot connect to %s database", DbDriver)
			log.Fatal("this is the error:", err)
		} else {
			fmt.Printf("we are connected to the %s database\n", DbDriver)
		}
	}

	s.DB.AutoMigrate(&models.User{}, &models.Post{})
	s.Router = mux.NewRouter()
	s.initializeRoutes()
}

func (s *Server) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
