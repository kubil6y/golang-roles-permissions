package data

import (
	"errors"

	"gorm.io/gorm"
)

type Role struct {
	CoreModel
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
}

type RoleModel struct {
	DB *gorm.DB
}

func (m RoleModel) GetAll() ([]*Role, error) {
	roles := make([]*Role, 0)
	err := m.DB.Find(&roles).Error
	if err != nil {
		return []*Role{}, err
	}
	return roles, nil
}

func (m RoleModel) GetByID(id int64) (*Role, error) {
	var role Role
	err := m.DB.First(&role, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}

func (m RoleModel) GetByName(name string) (*Role, error) {
	var role Role
	if err := m.DB.Where("name = ?", name).First(&role).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}

func (m RoleModel) Insert(r *Role) error {
	err := m.DB.Create(r).Error
	if err != nil {
		switch {
		// TODO
		case err.Error() == `ERROR: duplicate key value violates unique constraint "idx_users_email" (SQLSTATE 23505)`:
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

func (m RoleModel) Delete(r *Role) error {
	return m.DB.Delete(r).Error
}

func (m RoleModel) Update(r *Role) error {
	return m.DB.Model(r).Updates(r).Error
}
