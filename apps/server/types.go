package main

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                uuid.UUID `db:"id" json:"id"`
	Email             string    `db:"email" json:"email"`
	EncryptedPassword string    `json:"-"`
	Username          string    `db:"username" json:"username"`
	AvatarUrl         string    `db:"avatar_url" json:"avatar_url"`
	HeaderUrl         string    `db:"header_url" json:"header_url"`
	Bio               string    `db:"bio" json:"bio"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

func NewUser(email string, username string, password string, avatarUrl string, headerUrl string, bio string) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &User{
		Email:             email,
		EncryptedPassword: string(encryptedPassword),
		Username:          username,
		AvatarUrl:         avatarUrl,
		HeaderUrl:         headerUrl,
		Bio:               bio,
		CreatedAt:         time.Now().UTC(),
	}, nil
}
