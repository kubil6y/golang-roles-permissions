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

func (r *registerUserDTO) validate(v *validator.Validator) {
	v.Check(r.FirstName != "", "first_name", "must be provided")
	v.Check(r.LastName != "", "last_name", "must be provided")
	v.Check(r.Password != "", "password", "must be provided")
	v.Check(len(r.FirstName) > 1, "first_name", "must be longer than one character")
	v.Check(len(r.LastName) > 1, "last_name", "must be longer than one character")
	v.Check(len(r.Password) > 3, "password", "must be longer than three characters")

	validator.ValidateEmail(v, r.Email)
}

func (r *registerUserDTO) populate(user *data.User) {
	user.FirstName = r.FirstName
	user.LastName = r.LastName
	user.Email = r.Email
	user.SetPassword(r.Password) // hashing...
}
