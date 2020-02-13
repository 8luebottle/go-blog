package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//fmt.Println("WIP")

type User struct {
	ID        uint32    `gorm: "primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null; unique" json:"nickname"`
	Email     string    `gorm:"size:150;not null;" json:"email"`
	Password  string    `gorm:size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	ePw := "Please Input Your Password" // ePW == Empty Password
	eEm := "Please Input Your Email"    // eEm == Empty Email
	eNn := "Please Input Your Nickname" // eNn == Empty Nickname
	iEm := "Invalid Email"              // iEm == Invalid Email

	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New(eNn)
		}
		if u.Password == "" {
			return errors.New(ePw)
		}
		if u.Email == "" {
			return errors.New(eEm)
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New(iEm)
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("edit me")
		}
		if u.Email == "" {
			return errors.New(eEm)
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New(iEm)
		}
		return nil
	default:
		if u.Nickname == "" {
			return errors.New(eNn)
		}
		if u.Password == "" {
			return errors.New(ePw)
		}
		if u.Email == "" {
			return errors.New(eEm)
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New(iEm)
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error){
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil{
		return &[]User{}, err
	}
	return &usres, err
}