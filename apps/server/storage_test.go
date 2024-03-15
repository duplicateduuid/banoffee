package main

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func testRepositories(cleanup func(db *sqlx.DB), t *testing.T) Repositories {
	db, err := sqlx.Connect("postgres", "postgres://test:test@localhost:5433/test?sslmode=disable")
	if err != nil {
		fmt.Printf("[ERRO] %v", err)
	}

	userRepo := UserPostgresRepository{db: db}
	resourceRepo := ResourcePostgresRepository{db: db}
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	t.Cleanup(func() {
		cleanup(db)
	})

	return Repositories{
		userRepository:     userRepo,
		resourceRepository: resourceRepo,
		redis:              redis,
	}
}
