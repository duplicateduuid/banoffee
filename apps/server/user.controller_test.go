package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func deleteUserByEmail(email string) func(*sqlx.DB) {
	return func(db *sqlx.DB) {
		db.Exec(`DELETE FROM "user" WHERE email = $1`, email)
	}
}

func TestRegister(t *testing.T) {
	t.Parallel()

	repo := testRepositories(deleteUserByEmail("henri@henri.com"), t)
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
