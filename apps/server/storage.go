package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	bcrypt "golang.org/x/crypto/bcrypt"
)

type Repositories struct {
	userRepository UserRepository
}

type UserRepository interface {
	CreateUser(*User) error
	GetUser(uuid.UUID, User) error
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

func (u UserPostgresRepository) GetUser(id uuid.UUID, user User) error {
	return u.db.Get(user, `SELECT u.id, u.email FROM "user" u WHERE u.id=$1`, id)
}

func (u UserPostgresRepository) CreateUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), 10)

	if err != nil {
		return err
	}

	_, err = u.db.Exec(
		`INSERT INTO "user" (email, password, username)
		VALUES ($1, $2, $3)`,
		user.email,
		string(hashedPassword),
		user.username,
	)

	return err
}
