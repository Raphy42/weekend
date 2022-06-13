package users

import (
	"github.com/matthewhartstonge/argon2"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

var argon = argon2.DefaultConfig()

type User struct {
	gorm.Model
	Email        string `validate:"required,email"`
	PasswordHash []byte
}

func NewUser(email, password string) (*User, error) {
	passwordHash, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return nil, stacktrace.Propagate(err, "password hashing failed")
	}

	return &User{
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
