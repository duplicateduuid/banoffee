package main

import (
	"crypto/bcrypt"
	"github.com/jmoiron/sqlx"
	"hash"
)

type Repositories struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(*User) error
}

type UserPostgresRepository struct {
	db *sqlx.DB
}

func (u UserPostgresRepository) CreateUser(user *User) error {
	_, err := u.db.Exec(
		`INSERT INTO "user" (email, password, username, avatar_url, header_url, bio)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		user.email,
		user.username,
		user.avatarUrl,
		user.headerUrl,
		user.bio,
	)

	return err
}
