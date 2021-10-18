package data

import "gorm.io/gorm"

type Models struct {
	Users UserModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
