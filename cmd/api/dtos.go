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
	Name string `json:"name"`
}

func (d *roleDTO) validate(v *validator.Validator) {
	v.Check(d.Name != "", "name", "must be provided")
}

func (d *roleDTO) populate(r *data.Role) {
	r.Name = d.Name
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
