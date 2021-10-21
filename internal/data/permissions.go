package data

import (
	"errors"

	"gorm.io/gorm"
)

type Permission struct {
	CoreModel
	Name  string `json:"name" gorm:"uniqueIndex;not null"`
	Roles []Role `json:"permissions,omitempty" gorm:"many2many:roles_permissions;constraint:OnDelete:CASCADE"`
}

type PermissionModel struct {
	DB *gorm.DB
}

func (m PermissionModel) GetAll() ([]*Permission, error) {
	permissions := make([]*Permission, 0)
	err := m.DB.Find(&permissions).Error
	if err != nil {
		return []*Permission{}, err
	}
	return permissions, nil
}

func (m PermissionModel) GetByID(id int64) (*Permission, error) {
	var permission Permission
	err := m.DB.First(&permission, id).Error
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &permission, nil
}

func (m PermissionModel) GetByName(name string) (*Permission, error) {
	var permission Permission
	if err := m.DB.Where("name = ?", name).First(&permission).Error; err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &permission, nil
}

func (m PermissionModel) Insert(p *Permission) error {
	err := m.DB.Create(p).Error
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "idx_permissions_name" (SQLSTATE 23505)`:
			return ErrDuplicateRecord
		default:
			return err
		}
	}
	return nil
}

func (m PermissionModel) Delete(p *Permission) error {
	return m.DB.Delete(p).Error
}

func (m PermissionModel) Update(p *Permission) error {
	return m.DB.Model(p).Updates(p).Error
}
