package data

import (
	"errors"

	"gorm.io/gorm"
)

type Role struct {
	CoreModel
	Name        string       `json:"name" gorm:"uniqueIndex;not null"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:roles_permissions;constraint:OnDelete:CASCADE"`
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
	err := m.DB.Preload("Permissions").Where("id=?", id).First(&role).Error
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
	if err := m.DB.Preload("Permissions").Where("name = ?", name).First(&role).Error; err != nil {
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
		case IsDuplicateRecord(err):
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
	err := m.DB.Model(r).Association("Permissions").Replace(r.Permissions)
	err = m.DB.Save(r).Error
	return err
}
