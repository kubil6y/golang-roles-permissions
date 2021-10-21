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

func (p *PermissionModel) Count() int {
	var total int64
	p.DB.Model(&Permission{}).Count(&total)
	return int(total)
}

func (m PermissionModel) GetAll(p *Paginate) ([]*Permission, Metadata, error) {
	permissions := make([]*Permission, 0)
	err := m.DB.Scopes(p.PaginatedResults).Find(&permissions).Error
	if err != nil {
		return []*Permission{}, Metadata{}, err
	}

	var total int64
	m.DB.Model(&Permission{}).Count(&total)
	metadata := CalculateMetadata(p, int(total))

	return permissions, metadata, nil
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
		case IsDuplicateRecord(err):
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
