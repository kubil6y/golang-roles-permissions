package data

type Basket struct {
	CoreModel
	Elmalar []Elma `gorm:"foreignKey:BasketID"`
}
