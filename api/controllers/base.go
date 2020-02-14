package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	"github.com/jinzhu/gorm/dialects/mysql"

	"github.com/8luebottle/go-blog/api/models"
)

type Server struct {
	DB *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=ut8&parseTime=True&loc=Local",
			DbUser, DbPassword, DbHost, DbPort, DbName
			)
		if err != nil{
			fmt.Printf("cannot connect to %s database", Dbdriver)
			log.Fatal("this is the error:", err)
		}else {
			fmt.Printf("we are connected to the %s database", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

func (server *Server) Run(addr string){
	fmt.Println("listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}