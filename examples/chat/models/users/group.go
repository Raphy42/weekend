package users

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name    string `validate:"required"`
	Members []User
}
