package main

import (
	"fmt"
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
	CreateUser(*User) (*User, error)
	GetUserById(uuid.UUID) (*User, error)
	GetUserByEmail(string) (*User, error)
	GetUserByUsernameOrEmail(string) (*User, error)
	GetUserResource(user *User, resourceId string) (*Resource, error)
	GetUserResources(user *User, limit int, offset int, status string, reviewRating string) (*[]Resource, error)
	CreateUserResource(user *User, resourceId string, status *string, reviewRating *string, reviewComment *string) error
	UpdateUserResource(user *User, resourceId string, status *string, reviewRating *string, reviewComment *string) error
}

type UserPostgresRepository struct {
	db *sqlx.DB
}

type ResourceRepository interface {
	CreateResource(*Resource) (*Resource, error)
	GetResourceById(string) (*Resource, error)
	GetResourceByUrl(string) (*Resource, error)
	GetResourceByName(name string) (*Resource, error)
	GetUserResources(user *User, limit int, offset int) (*[]Resource, error)
	SearchByName(name string, limit int, offset int) (*[]Resource, error)
}

type ResourcePostgresRepository struct {
	db *sqlx.DB
}

func NewRepositories() Repositories {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=banoffee password=5up3r_s3cur3_p4ssw0rd sslmode=disable")

	if err != nil {
		log.Fatalln(err)
	}

	userRepo := UserPostgresRepository{db: db}
	resourceRepo := ResourcePostgresRepository{db: db}
	// TODO: this is bad. write an actual redis repository
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return Repositories{
		userRepository:     userRepo,
		resourceRepository: resourceRepo,
		redis:              redis,
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

func (u UserPostgresRepository) GetUserByUsernameOrEmail(usernameOrEmail string) (*User, error) {
	user := new(User)
	err := u.db.Get(
		user,
		`SELECT u.id, u.email, u.username, u.avatar_url, u.header_url, u.bio
		 FROM "user" u 
		 WHERE u.email=$1 OR u.username=$1`,
		usernameOrEmail,
	)

	return user, err
}

func (u UserPostgresRepository) CreateUser(user *User) (*User, error) {
	_, err := u.db.Exec(
		`INSERT INTO "user" (email, password, username)
		VALUES ($1, $2, $3)`,
		user.Email,
		user.EncryptedPassword,
		user.Username,
	)

	if err != nil {
		return nil, err
	}

	return u.GetUserByEmail(user.Email)
}

func (u UserPostgresRepository) GetUserResource(user *User, resourceId string) (*Resource, error) {
	resource := new(Resource)

	err := u.db.Get(
		resource,
		`
		SELECT 
			r.id, r.url, r.name, r.image_url, r.author, r.description, ur.status, ur.review_rating, ur.review_comment, r.created_at
		FROM 
			"user_resource" ur 
		LEFT JOIN 
			"resource" r ON r.id = ur.resource_id
		WHERE
			ur.user_id = $1
			AND ur.resource_id = $2
		`,
		user.Id,
		resourceId,
	)

	return resource, err
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

func (r ResourcePostgresRepository) CreateResource(resource *Resource) (*Resource, error) {
	_, err := r.db.Exec(
		`INSERT INTO "resource" (url, name, image_url, author, description)
		VALUES ($1, $2, $3, $4, $5)`,
		resource.Url,
		resource.Name,
		resource.ImageUrl,
		resource.Author,
		resource.Description,
	)

	if err != nil {
		return nil, err
	}

	return r.GetResourceByUrl(resource.Url)
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

func (r ResourcePostgresRepository) SearchByName(name string, limit int, offset int) (*[]Resource, error) {
	resources := new([]Resource)
	name = "%" + name + "%"

	err := r.db.Select(
		resources,
		`
		SELECT 
		  r.id, r.url, r.name, r.image_url, r.author, r.description, r.created_at
		FROM 
		  "resource" r
		ORDER BY 
			SIMILARITY(r.name, $1) DESC
		LIMIT
			$2
		OFFSET
			$3
		`,
		name,
		limit,
		offset,
	)

	fmt.Printf("[INFO] [ResourcePostgresRepository.SearchByName] selected: %v\n", resources)

	return resources, err
}
