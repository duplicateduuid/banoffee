package main

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testNewUser(t *testing.T, repos *Repositories) (*User, string) {
	username := uuid.New().String()
	email := username + "@x.com"
	password := "testestest"

	user, err := NewUser(email, username, password, nil, nil, nil)
	user, err = repos.userRepository.CreateUser(user)

	if err != nil {
		t.Errorf("[ERROR] failed to create auth user: %s\n", err)
	}

	return user, password
}

func TestLogin(t *testing.T) {
	t.Parallel()

	repos := newTestRepositories(t)
	user, password := testNewUser(t, &repos)

	req := LoginRequest{
		Email:    user.Email,
		Password: password,
	}
	router := newTestRouter(t, repos)
	w := router.post("/login", req)

	body := w.Body.String()
	assert.Equal(t, 200, w.Code, body)

	var resp LoginResponse
	json.Unmarshal([]byte(body), &resp)
	assert.Equal(t, user.Email, resp.User.Email)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	repos := newTestRepositories(t)
	router := newTestRouter(t, repos)

	username := uuid.New().String()[:10]
	password := uuid.New().String()

	req := RegisterRequest{
		Email:    username + "@henri.com",
		Password: password,
		Username: username,
	}
	w := router.post("/register", req)

	body := w.Body.String()
	assert.Equal(t, 200, w.Code, body)

	var response RegisterResponse
	json.Unmarshal([]byte(body), &response)

	assert.Equal(t, req.Email, response.User.Email)
}
