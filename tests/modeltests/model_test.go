package modeltests

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
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TestDBDriver")

	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charse=utf8&parseTime=True&loc=Local",
			os.Getenv("TestDbUser"),
			os.Getenv("TestDbPassword"),
			os.Getenv("TestDbHost"),
			os.Getenv("TestDbPort"),
			os.Getenv("TestDbName"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("cannot connect to %s database\n", TestDbDriver)
		} else {
			fmt.Printf("we are conntected to the %s database\n", TestDbDriver)
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
	refreshUserTable()

	user := models.User{
		Nickname: "testuser",
		Email:    "testuser@gmail.com",
		Password: "test123",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {
	users := []models.User{
		models.User{
			Nickname: "Henry",
			Email:    "henrylau@gmial.com",
			Password: "henry123",
		},
		models.User{
			Nickname: "Trump",
			Email:    "donaldtrump@gmail.com",
			Password: "donald123",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
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
		Nickname: "Elon",
		Email:    "elonmusk@gmail.com",
		Password: "elon123",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		Title:    "Space X",
		Content:  "It's time for another flying Falcon 9 to deliver a fifth batch of Starlink satellites to orbit for SpaceX. On Sunday morning, SpaceX's workhorse rocket is set to add 60 more satellites to the burgeoning constellation.",
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
			Nickname: "Elon",
			Email:    "elonmusk@gmail.com",
			Password: "elon123",
		},
		models.User{
			Nickname: "Adult Tiger",
			Email:    "tiger@gmail.com",
			Password: "tiger123",
		},
	}
	var posts = []models.Post{
		models.Post{
			Title:   "How to watch Falcon 9 deliver 60 satellites to space",
			Content: "SpaceX attempts to go five for five with its Starlink satellite megaconstellation.",
		},
		models.Post{
			Title:   "You Can Now Lease A Tesla In Connecticut",
			Content: "The new leasing option was a cause for celebration among Tesla officials, EV advocates, and state officials despite the fact that state legislators haven’t changed the state’s law requiring third-party dealerships for consumer sales. Dealerships have argued that this would be unfair and that franchising helps drive competition and sales.",
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
