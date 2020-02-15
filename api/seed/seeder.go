package seed

import (
	"log"

	"github.com/8luebottle/go-blog/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: "Baby Tiger",
		Email:    "babytiger@gmal.com",
		Password: "baby123",
	},
	models.User{
		Nickname: "Thom Browne",
		Email:    "thombrowne@gmail.com",
		Password: "thom123",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "First Post",
		Content: "Go Programming. Daily Coding",
	},
	models.Post{
		Title:   "Second Post",
		Content: "Level Up Coding.",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
