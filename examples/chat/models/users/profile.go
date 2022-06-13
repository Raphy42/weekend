package users

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Username      string `validate:"required"`
	UsernameColor string `validate:"required,hsl"`
}
