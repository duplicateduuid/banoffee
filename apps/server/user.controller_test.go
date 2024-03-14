package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func repo(t *testing.T) Repositories {
	db, err := sqlx.Connect("postgres", "postgres://test:test@localhost:5433/test?sslmode=disable")
	if err != nil {
		fmt.Printf("[ERRO] %v", err)
	}

	userRepo := UserPostgresRepository{
		db: db,
	}

	t.Cleanup(func() {
		deleteUserByEmail("henri@henri.com", db)
	})

	return Repositories{
		userRepository: userRepo,
	}
}

func deleteUserByEmail(email string, db *sqlx.DB) {
	db.Exec(`DELETE FROM "user" WHERE email = $1`, email)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	repo := repo(t)
	router := NewAPI(repo).SetupRouter()

	user := RegisterRequest{
		Email:    "henri@henri.com",
		Password: "henrihenrihenri",
		Username: "henri",
	}
	json_user, _ := json.Marshal(user)
	payload := bytes.NewReader(json_user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", payload)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"User created with success!"}`, w.Body.String())
}
