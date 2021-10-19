package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type CoreModel struct {
	ID        int            `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Models struct {
	Users UserModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users: UserModel{DB: db},
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
