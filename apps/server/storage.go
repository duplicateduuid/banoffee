package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(*User) error
	GetUserById(uuid.UUID, User) error
	GetUserByEmail(string, User) error
}

type UserPostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepositories() *Repositories {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=banoffee password=5up3r_s3cur3_p4ssw0rd sslmode=disable")

	if err != nil {
		log.Fatalln(err)
	}

	userRepo := UserPostgresRepository{
		db: db,
	}

	return &Repositories{
		userRepository: userRepo,
	}
}

type UserAuthInfo struct {
	Id    uuid.UUID `db:"id" json:"id"`
	Email string    `db:"email" json:"email"`
}

func (u UserPostgresRepository) GetUserById(id uuid.UUID, user User) error {
	return u.db.Get(
		user,
		`SELECT u.id, u.email u.password u.username u.avatar_url u.header_url u.bio
		FROM "user" u WHERE u.id=$1`,
		id,
	)
}

func (u UserPostgresRepository) GetUserByEmail(email string, user User) error {
	return u.db.Get(
		user,
		`SELECT u.id, u.email u.password u.username u.avatar_url u.header_url u.bio
		FROM "user" u WHERE u.email=$1`,
		email,
	)
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
