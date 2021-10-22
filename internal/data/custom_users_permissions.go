package data

import "gorm.io/gorm"

type CustomUsersPermissions struct {
	CoreModel
	UserID       int64      `json:"user_id" gorm:"not null"`
	User         User       `json:"user,omitempty"`
	PermissionID int64      `json:"permission_id" gorm:"not null"`
	Permission   Permission `json:"permission,omitempty"`
	HasAccess    bool       `json:"has_access" gorm:"not null"`
}

type CustomUsersPermissionsModel struct {
	DB *gorm.DB
}

func (m CustomUsersPermissionsModel) Insert(c *CustomUsersPermissions) error {
	err := m.DB.Create(c).Error
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

func (m CustomUsersPermissionsModel) Delete(c *CustomUsersPermissions) error {
	return m.DB.Delete(c).Error
}

func (m CustomUsersPermissionsModel) Update(c *CustomUsersPermissions) error {
	return m.DB.Updates(c).Error
}
