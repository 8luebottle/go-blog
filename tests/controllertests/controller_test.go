package controllertests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/8luebottle/go-blog/api/controllers"
	"github.com/8luebottle/go-blog/api/models"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var postInstance = models.Post{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("error getting env %v\n", err)
	}
	Database()
	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TestDbDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			os.Getenv("TestDbUser"),
			os.Getenv("TestDbPassword"),
			os.Getenv("TestDbHost"),
			os.Getenv("TestDbPort"),
			os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("cannot connect to %s database\n", TestDbDriver)
			log.Fatal("this is the error:", err)
		} else {
			fmt.Printf("we are connected to the %s database\n", TestDbDriver)
		}
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Nickname: "Test",
		Email:    "test@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {
	var err error
	if err != nil {
		return nil, err
	}
	users := []models.User{
		models.User{
			Nickname: "Jeff",
			Email:    "jeff@gmail.com",
			Password: "jeff123",
		},
		models.User{
			Nickname: "Tonny",
			Email:    "morry@gmail.com",
			Password: "tonny123",
		},
	}
	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndPostTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("successfully refreshed tables")
	return nil
}

func seedOneUserAndOnePost() (models.Post, error) {
	err := refreshUserAndPostTable()
	if err != nil {
		return models.Post{}, err
	}
	user := models.User{
		Nickname: "Madison Beer",
		Email:    "madison@gmail.com",
		Password: "madison123",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		Title:    "Good In Goodbye",
		Content:  "You got blood on your hands, how do you plead?",
		AuthorID: user.ID,
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func seedUsersAndPosts() ([]models.User, []models.Post, error) {
	var err error
	if err != nil {
		return []models.User{}, []models.Post{}, err
	}
	var users = []models.User{
		models.User{
			Nickname: "Itzy",
			Email:    "ityzy@gmail.com",
			Password: "itzy123",
		},
		models.User{
			Nickname: "Twice",
			Email:    "twice@gmail.com",
			Password: "twice123",
		},
	}
	var posts = []models.Post{
		models.Post{
			Title:   "ITZY's New song",
			Content: "https://www.youtube.com/watch?v=W8yt_avdKbY",
		},
		models.Post{
			Title:   "Twice's New Song",
			Content: "NONE",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}
