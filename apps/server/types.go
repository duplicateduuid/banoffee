package main

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	id         uuid.UUID `json:"id"`
	email      string    `json:"email"`
	password   string    `json:"-`
	username   string    `json:"username"`
	avatarUrl  string    `json:"avatar_url"`
	headerUrl  string    `json:"header_url"`
	bio        string    `json:"bio"`
	created_at time.Time `json:"created_at"`
}
