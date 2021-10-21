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
	FirstName   string       `json:"first_name" gorm:"not null"`
	LastName    string       `json:"last_name" gorm:"not null"`
	Email       string       `json:"email" gorm:"uniqueIndex;not null"`
	Password    []byte       `json:"-" gorm:"not null"`
	IsActivated bool         `json:"-" gorm:"default:false;not null"`
	IsAdmin     bool         `json:"-" gorm:"default:false;not null"`
	Tokens      []Token      `json:"tokens,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:users_permissions"`
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
		case IsDuplicateRecord(err):
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) Update(u *User) error {
	return m.DB.Model(u).Updates(u).Error
}

func (m UserModel) Delete(u *User) error {
	return m.DB.Model(u).Delete(u).Error
}

func (m UserModel) GetByID(id int64) (*User, error) {
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
	err := m.DB.Where("hash=? and scope=? and expiry > ?", tokenHash, scope, time.Now()).Preload("User").First(&token).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &token.User, nil
}

func (m UserModel) GetAll(p *Paginate) ([]*User, Metadata, error) {
	users := make([]*User, 0)
	err := m.DB.Scopes(p.PaginatedResults).Find(&users).Error
	if err != nil {
		return nil, Metadata{}, nil
	}

	var total int64
	m.DB.Model(&User{}).Count(&total)
	metadata := CalculateMetadata(p, int(total))
	return users, metadata, nil
}
