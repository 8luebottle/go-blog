package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:150;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
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
	ePw := "please input your password" // ePW == Empty Password
	eEm := "please input your email"    // eEm == Empty Email
	eNn := "please input your nickname" // eNn == Empty Nickname
	iEm := "invalid email"              // iEm == Invalid Email

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
			return errors.New(ePw)
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
	if err := db.Create(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	if err := db.Model(&User{}).
		Limit(100).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	err := db.Model(User{}).
		Where("id = ?", uid).
		Take(&u).Error
	if err != nil {
		return nil, err
	}

	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Model(&User{}).
		Where("id = ?", uid).
		Take(&User{}).
		UpdateColumns(
			map[string]interface{}{
				"password":  u.Password,
				"nickname":  u.Nickname,
				"email":     u.Email,
				"update_at": time.Now(),
			},
		)

	if db.Error != nil {
		return nil, db.Error
	}

	if err = db.Model(&User{}).
		Where("id = ?", uid).
		Take(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Model(&User{}).
		Where("id = ?", uid).
		Take(&User{}).
		Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
