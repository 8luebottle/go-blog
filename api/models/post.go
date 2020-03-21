package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Title     string     `gorm:"size:255;not null;unique" json:"title"`
	Content   string     `gorm:"size:500;not null;" json:"content"`
	Author    User       `json:"author"`
	AuthorID  uint32     `gorm:"not null" json:"author_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (p *Post) Prepare() {
	now := time.Now()
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = now
	p.UpdatedAt = now
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New("required Title")
	}
	if p.Content == "" {
		return errors.New("required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("required Author")
	}
	return nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	if err := db.Model(&Post{}).
		Create(&p).Error; err != nil {
		return nil, err
	}

	if p.ID != 0 {
		if err := db.Model(&User{}).
			Where("id = ?", p.AuthorID).
			Take(&p.Author).Error; err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	posts := []Post{}
	if err := db.Model(&Post{}).
		Limit(100).
		Find(&posts).Error; err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.
				Model(&User{}).
				Where("id = ?", posts[i].AuthorID).
				Take(&posts[i].Author).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return &posts, nil
}

func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	if err := db.Model(&Post{}).
		Where("id = ?", pid).
		Take(&p).Error; err != nil {
		return nil, err
	}

	if p.ID != 0 {
		if err := db.Model(&User{}).
			Where("id = ?", p.AuthorID).
			Take(&p.Author).Error; err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {
	if err := db.Model(&Post{}).
		Where("id = ?", p.ID).
		Update(
			Post{
				Title:     p.Title,
				Content:   p.Content,
				UpdatedAt: time.Now(),
			}).Error; err != nil {
		return nil, err
	}

	if p.ID != 0 {
		if err := db.Model(&User{}).
			Where("id = ?", p.AuthorID).
			Take(&p.Author).Error; err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Post) DeleteAPost(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Model(&Post{}).
		Where("id = ? and author_id = ?", pid, uid).
		Take(&Post{}).
		Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
