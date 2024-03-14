package main

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	userRepository     UserRepository
	resourceRepository ResourceRepository
}

type UserRepository interface {
	CreateUser(*User) error
	GetUserById(uuid.UUID) (*User, error)
	GetUserByEmail(string) (*User, error)
}

type UserPostgresRepository struct {
	db *sqlx.DB
}

type ResourceRepository interface {
	CreateResource(*Resource) error
	GetResourceById(uuid.UUID) (*Resource, error)
	GetResourceByUrl(string) (*Resource, error)
}

type ResourcePostgresRepository struct {
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

func (u UserPostgresRepository) GetUserById(id uuid.UUID) (*User, error) {
	user := new(User)
	err := u.db.Get(
		user,
		`SELECT u.id, u.email, u.password, u.username, u.avatar_url, u.header_url, u.bio
		FROM "user" u WHERE u.id=$1`,
		id,
	)

	return user, err
}

func (u UserPostgresRepository) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	err := u.db.Get(
		user,
		`SELECT u.id, u.email, u.password, u.username, u.avatar_url, u.header_url, u.bio
		FROM "user" u WHERE u.email=$1`,
		email,
	)

	return user, err
}

func (u UserPostgresRepository) CreateUser(user *User) error {
	_, err := u.db.Exec(
		`INSERT INTO "user" (email, password, username)
		VALUES ($1, $2, $3)`,
		user.Email,
		user.EncryptedPassword,
		user.Username,
	)

	return err
}

func (r ResourcePostgresRepository) CreateResource(resource *Resource) error {
	imageUrl := sql.NullString{String: resource.ImageUrl, Valid: resource.ImageUrl != ""}
	author := sql.NullString{String: resource.Author, Valid: resource.Author != ""}
	description := sql.NullString{String: resource.Description, Valid: resource.Description != ""}

	_, err := r.db.Exec(
		`INSERT INTO "resource" (url, name, image_url, author, description)
		VALUES ($1, $2, $3, $4, $5)`,
		resource.Url,
		resource.Name,
		imageUrl,
		author,
		description,
	)

	return err
}

func (r ResourcePostgresRepository) GetResourceById(id uuid.UUID) (*Resource, error) {
	resource := new(Resource)

	err := r.db.Get(
		resource,
		`SELECT r.id, r.url, r.name, r.image_url, r.author, r.description, r.created_at FROM "resource" r WHERE r.id=$1`,
		id,
	)

	return resource, err
}

func (r ResourcePostgresRepository) GetResourceByUrl(url string) (*Resource, error) {
	resource := new(Resource)

	err := r.db.Get(
		resource,
		`SELECT r.id, r.url, r.name, r.image_url, r.author, r.description, r.created_at FROM "resource" r WHERE r.url=$1`,
		url,
	)

	return resource, err
}
