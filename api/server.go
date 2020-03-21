package api

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/8luebottle/go-blog/api/controllers"
	"github.com/8luebottle/go-blog/api/seed"
)

var server = controllers.Server{}

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error getting env, not comming through %v", err)
	}

	server.Initialize(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	seed.Load(server.DB)
	server.Run(":8000")
}
