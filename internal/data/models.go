package data

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type CoreModel struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Models struct {
	Users       UserModel
	Tokens      TokenModel
	Roles       RoleModel
	Permissions PermissionModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Roles:       RoleModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}

func IsDuplicateRecord(err error) bool {
	fmt.Println("this is hit")
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.Code == "23505"
		}
	}
	return false
}
