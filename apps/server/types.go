package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id         uuid.UUID `db:"id" json:"id"`
	Email      string    `db:"email" json:"email"`
	Password   string    `json:"password"`
	Username   string    `db:"username" json:"username"`
	AvatarUrl  *string   `db:"avatar_url" json:"avatar_url"`
	HeaderUrl  *string   `db:"header_url" json:"header_url"`
	Bio        *string   `db:"bio" json:"bio"`
	Created_at time.Time `json:"created_at"`
}
