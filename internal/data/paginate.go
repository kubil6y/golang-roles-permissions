// inspired by this source + alex: https://medium.com/@michalkowal567/creating-reusable-pagination-in-golang-and-gorm-4b23e179a54b
package data

import (
	"math"

	"github.com/kubil6y/myshop-go/internal/validator"
	"gorm.io/gorm"
)

type Paginate struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

// PaginatedResults is used when making db calls, example:
// err := m.DB.Scopes(p.PaginatedResults).Find(&permissions).Error
func (p Paginate) PaginatedResults(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return db.Offset(offset).Limit(p.Limit)
}

func ValidatePaginate(v *validator.Validator, p *Paginate) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Limit > 0, "limit", "must be greater than zero")
	v.Check(p.Page <= 10_000, "page", "must be a maximum of 10_000")
	v.Check(p.Limit <= 100, "limit", "must be a maximum of 100")
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func CalculateMetadata(p *Paginate, total int) Metadata {
	if total == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  p.Page,
		PageSize:     p.Limit,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(total) / float64(p.Limit))),
		TotalRecords: total,
	}
}
