package data

import "gorm.io/gorm"

type Elma struct {
	gorm.Model
	BasketID int64
}
