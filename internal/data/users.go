package data

import (
	"crypto/sha256"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	AnonymousUser = &User{}
)

type User struct {
	CoreModel
	FirstName   string  `json:"first_name" gorm:"not null"`
	LastName    string  `json:"last_name" gorm:"not null"`
	Email       string  `json:"email" gorm:"uniqueIndex;not null"`
	Password    []byte  `json:"-" gorm:"not null"`
	IsActivated bool    `json:"-"`
	Tokens      []Token `json:"tokens" gorm:"foreignKey:UserID"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) SetPassword(plain string) error {
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}

func (u *User) ComparePassword(plain string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(plain))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type UserModel struct {
	DB *gorm.DB
}

func (m UserModel) Insert(u *User) error {
	if err := m.DB.Create(u).Error; err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "idx_users_email" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) Update(u *User) error {
	return m.DB.Model(u).Updates(u).Error
}

func (m UserModel) GetById(id int64) (*User, error) {
	var user User
	if err := m.DB.First(&user, id).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	var user User
	if err := m.DB.Where("email = ?", email).First(&user).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) GetForToken(scope string, tokenPlaintext string) (*User, error) {
	sizedTokenHash := sha256.Sum256([]byte(tokenPlaintext))
	tokenHash := sizedTokenHash[:]

	var token Token
	err := m.DB.Where("hash=? and scope=? and expiry > ?", tokenHash, scope, time.Now()).Preload("User").Find(&token).Error
	if err != nil {
		return nil, err
	}

	if token.User.ID == 0 {
		return nil, ErrRecordNotFound
	}

	return &token.User, nil
}