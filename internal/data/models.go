package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type CoreModel struct {
	ID        int64          `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Models struct {
	Users  UserModel
	Tokens TokenModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}

func IsDuplicateRecord(err error) {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message)
			fmt.Println(pgErr.Code)
		}
	}
}
