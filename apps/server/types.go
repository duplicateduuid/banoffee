package main

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                uuid.UUID `db:"id" json:"id" tstype:"string"`
	Email             string    `db:"email" json:"email" tstype:"string"`
	EncryptedPassword string    `db:"password" json:"-" tstype:"string"`
	Username          string    `db:"username" json:"username" tstype:"string"`
	AvatarUrl         *string   `db:"avatar_url" json:"avatar_url" tstype:"string | null"`
	HeaderUrl         *string   `db:"header_url" json:"header_url" tstype:"string | null"`
	Bio               *string   `db:"bio" json:"bio" tstype:"string | null"`
	CreatedAt         time.Time `db:"created_at" json:"created_at" tstype:"string"`
}

type Resource struct {
	Id            uuid.UUID `db:"id" json:"id" tstype:"string"`
	Url           string    `db:"url" json:"url" tstype:"string"`
	Name          string    `db:"name" json:"name" tstype:"string"`
	ImageUrl      *string   `db:"image_url" json:"image_url" tstype:"null | string"`
	Author        *string   `db:"author" json:"author" tstype:"string | null"`
	Description   *string   `db:"description" json:"description" tstype:"string | null"`
	Status        *string   `db:"status" json:"status" tstype:"string | null"`
	ReviewRating  *string   `db:"review_rating" json:"review_rating" tstype:"string | null"`
	ReviewComment *string   `db:"review_comment" json:"review_comment" tstype:"string | null"`
	CreatedAt     time.Time `db:"created_at" json:"created_at" tstype:"string"`
}

func (u *User) ValidPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) != nil
}

func NewUser(email string, username string, password string, avatarUrl *string, headerUrl *string, bio *string) (*User, error) {
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

func NewResource(url string, name string, imageUrl *string, author *string, description *string) *Resource {
	return &Resource{
		Url:         url,
		Name:        name,
		ImageUrl:    imageUrl,
		Author:      author,
		Description: description,
	}
}
