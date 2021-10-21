package main

import (
	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

type registerUserDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (d *registerUserDTO) validate(v *validator.Validator) {
	v.Check(d.FirstName != "", "first_name", "must be provided")
	v.Check(d.LastName != "", "last_name", "must be provided")
	v.Check(d.Password != "", "password", "must be provided")
	v.Check(len(d.FirstName) > 1, "first_name", "must be longer than one character")
	v.Check(len(d.LastName) > 1, "last_name", "must be longer than one character")
	v.Check(len(d.Password) > 3, "password", "must be longer than three characters")

	validator.ValidateEmail(v, d.Email)
}

func (d *registerUserDTO) populate(user *data.User) {
	user.FirstName = d.FirstName
	user.LastName = d.LastName
	user.Email = d.Email
}

type createAuthenticationTokenDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *createAuthenticationTokenDto) validate(v *validator.Validator) {
	validator.ValidateEmail(v, d.Email)
	v.Check(d.Password != "", "password", "must be provided")
	v.Check(len(d.Password) > 3, "password", "must be longer than three characters")
}

type permissionDTO struct {
	Name string `json:"name"`
}

func (d *permissionDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
}
func (d *permissionDTO) populate(p *data.Permission) {
	p.Name = d.Name
}

type roleDTO struct {
	Name        string  `json:"name"`
	Permissions []int64 `json:"permissions"`
}

func (d *roleDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
	v.Check(len(d.Permissions) != 0, "permissions", "must be provided")
	v.Check(validator.IsUniqueIS(d.Permissions), "permissions", "values must be unique")
}

type grantPermissionsToRolesDTO struct {
	RoleID      int64   `json:"role_id"`
	Permissions []int64 `json:"permissions"`
}

func (d *grantPermissionsToRolesDTO) validate(v *validator.Validator) {
	v.Check(d.RoleID > 0, "role_id", "invalid value")
	v.Check(len(d.Permissions) > 0, "permissions", "invalid value, []int > 0")
	v.Check(validator.IsUniqueIS(d.Permissions), "permissions", "must be unique values")
}

type updateUserDTO struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}

func (d *updateUserDTO) validate(v *validator.Validator) {
	if d.FirstName != nil {
		v.Check(*d.FirstName != "", "first_name", "can not be empty")
		v.Check(len(*d.FirstName) > 1, "first_name", "must be longer than one character")
	}

	if d.LastName != nil {
		v.Check(*d.LastName != "", "last_name", "can not be empty")
		v.Check(len(*d.LastName) > 1, "last_name", "must be longer than one character")
	}

	if d.Email != nil {
		v.Check(*d.Email != "", "email", "can not be empty")
		validator.ValidateEmail(v, *d.Email)
	}

	if d.Password != nil {
		v.Check(*d.Password != "", "password", "can not be empty")
		v.Check(len(*d.Password) > 3, "password", "must be longer than three characters")
	}
}

func (d *updateUserDTO) populate(u *data.User) error {
	var err error
	if d.FirstName != nil {
		u.FirstName = *d.FirstName
	}

	if d.LastName != nil {
		u.LastName = *d.LastName
	}

	if d.Email != nil {
		u.Email = *d.Email
	}

	if d.Password != nil {
		err = u.SetPassword(*d.Password)
	}
	return err
}
