package main

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	userRepository     UserRepository
	resourceRepository ResourceRepository
	// TODO: replace by an actual repository
	redis *redis.Client
}

type UserRepository interface {
	CreateUser(*User) error
	GetUserById(uuid.UUID) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserResources(user *User, limit int, offset int, status string, reviewRating string) (*[]Resource, error)
	GetUserResource(user *User, resourceId string) (*Resource, error)
	CreateUserResource(user *User, resourceId string, status *string, reviewRating *string, reviewComment *string) error
	UpdateUserResource(user *User, resourceId string, status *string, reviewRating *string, reviewComment *string) error
}

type UserPostgresRepository struct {
	db *sqlx.DB
}

type ResourceRepository interface {
	CreateResource(*Resource) error
	GetResourceById(string) (*Resource, error)
	GetResourceByUrl(string) (*Resource, error)
	GetResourceByName(name string) (*Resource, error)
	GetUserResources(user *User, limit int, offset int) (*[]Resource, error)
	SearchByName(name string) (*[]Resource, error)
}

type ResourcePostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepositories() Repositories {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=banoffee password=5up3r_s3cur3_p4ssw0rd sslmode=disable")

	if err != nil {
		log.Fatalln(err)
	}

	userRepo := UserPostgresRepository{db: db}
	resourceRepo := ResourcePostgresRepository{db: db}

	return Repositories{
		userRepository:     userRepo,
		resourceRepository: resourceRepo,
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
		`SELECT u.id, u.email, u.username, u.avatar_url, u.header_url, u.bio
		FROM "user" u WHERE u.id=$1`,
		id,
	)

	return user, err
}

func (u UserPostgresRepository) GetUserByEmail(email string) (*User, error) {
	user := new(User)
	err := u.db.Get(
		user,
		`SELECT u.id, u.email, u.username, u.avatar_url, u.header_url, u.bio
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

func (u UserPostgresRepository) GetUserResources(user *User, limit int, offset int, status string, reviewRating string) (*[]Resource, error) {
	resource := new([]Resource)

	err := u.db.Select(
		resource,
		`
		SELECT 
			r.id, r.url, r.name, r.image_url, r.author, r.description, ur.status, ur.review_rating, ur.review_comment, r.created_at
		FROM 
			"resource" r 
		LEFT JOIN 
			"user_resource" ur ON ur.resource_id = r.id
		WHERE
			ur.user_id = $1
			AND ($4 = '' OR ur.status::text = $4)
            AND ($5 = '' OR ur.review_rating::text = $5)
		ORDER BY
			r.created_at
		LIMIT
			$2
		OFFSET
			$3
		`,
		user.Id,
		limit,
		offset,
		status,
		reviewRating,
	)

	return resource, err
}

func (u UserPostgresRepository) GetUserResource(user *User, resourceId string) (*Resource, error) {
	resource := new(Resource)

	err := u.db.Get(
		resource,
		`
		SELECT 
			r.id, r.url, r.name, r.image_url, r.author, r.description, ur.status, ur.review_rating, ur.review_comment, r.created_at
		FROM 
			"resource" r 
		LEFT JOIN 
			"user_resource" ur ON ur.resource_id = r.id
		WHERE
			ur.user_id = $1
			AND r.id = $2
		`,
		user.Id,
		resourceId,
	)

	return resource, err
}

func (u UserPostgresRepository) CreateUserResource(user *User, resourceId string, status *string, reviewRating *string, reviewComment *string) error {
	newStatus := status

	if reviewRating != nil && status == nil {
		updatedStatus := "ongoing"
		newStatus = &updatedStatus
	}

	_, err := u.db.Exec(
		`INSERT INTO "user_resource" (user_id, resource_id, status, review_rating, review_comment)
		VALUES ($1, $2, $3, $4, $5)`,
		user.Id,
		resourceId,
		newStatus,
		reviewRating,
		reviewComment,
	)

	return err
}

func (u UserPostgresRepository) UpdateUserResource(user *User, resourceId string, status *string, reviewRating *string, reviewComment *string) error {
	_, err := u.db.Exec(
		`
			UPDATE "user_resource"
			SET
				status = COALESCE($3, status),
				review_rating = COALESCE($4, review_rating),
				review_comment = COALESCE($5, review_comment)
			WHERE
				user_id = $1
				AND resource_id = $2
		`,
		user.Id,
		resourceId,
		status,
		reviewRating,
		reviewComment,
	)

	return err
}

func (r ResourcePostgresRepository) CreateResource(resource *Resource) error {
	_, err := r.db.Exec(
		`INSERT INTO "resource" (url, name, image_url, author, description)
		VALUES ($1, $2, $3, $4, $5)`,
		resource.Url,
		resource.Name,
		resource.ImageUrl,
		resource.Author,
		resource.Description,
	)

	return err
}

func (r ResourcePostgresRepository) GetResourceById(id string) (*Resource, error) {
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

func (r ResourcePostgresRepository) GetResourceByName(name string) (*Resource, error) {
	resource := new(Resource)

	err := r.db.Get(
		resource,
		`SELECT r.id, r.url, r.name, r.image_url, r.author, r.description, r.created_at FROM "resource" r WHERE r.name=$1`,
		name,
	)

	return resource, err
}

func (r ResourcePostgresRepository) GetUserResources(user *User, limit int, offset int) (*[]Resource, error) {
	resource := new([]Resource)

	err := r.db.Select(
		resource,
		`
		SELECT 
			r.id, r.url, r.name, r.image_url, r.author, r.description, r.created_at
		FROM 
			"resource" r 
		LEFT JOIN 
			"user_resource" ur ON ur.resource_id = r.id
		WHERE
			ur.user_id = $1
		ORDER BY
			r.created_at
		LIMIT
			$2
		OFFSET
			$3
		`,
		user.Id,
		limit,
		offset,
	)

	return resource, err
}

func (r ResourcePostgresRepository) SearchByName(name string) (*[]Resource, error) {
	var resources []Resource

	err := r.db.Select(
		resources,
		`
		SELECT 
			r.id, r.url, r.name, r.image_url, r.author, r.description, r.created_at
		FROM 
			"resource" r 
		LEFT JOIN 
			"user_resource" ur ON ur.resource_id = r.id
		WHERE
			ur.user_id = $1
		ORDER BY
			r.created_at
		LIMIT
			$2
		OFFSET
			$3
		`,
		name,
	)

	return &resources, err
}
