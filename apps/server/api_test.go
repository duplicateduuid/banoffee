package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func newTestRepositories(cleanup func(db *sqlx.DB), t *testing.T) Repositories {
	db, err := sqlx.Connect("postgres", "postgres://test:test@localhost:5433/test?sslmode=disable")
	if err != nil {
		fmt.Printf("[ERRO] %v", err)
	}

	userRepo := UserPostgresRepository{db: db}
	resourceRepo := ResourcePostgresRepository{db: db}
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
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

type testRouter struct {
	router *gin.Engine
	repos  *Repositories
	auth   *User
}

func newTestRouter(repos Repositories) testRouter {
	router := NewAPI(repos).SetupRouter()

	return testRouter{router: router}
}

func newAuthTestRouter(repos Repositories, auth *User) testRouter {
	router := NewAPI(repos).SetupRouter()

	if auth == nil {
		user, _ := NewUser("test@test.com", "test", "testtesttest", nil, nil, nil)
		repos.userRepository.CreateUser(user)
		auth, _ = repos.userRepository.GetUserByEmail(user.Email)
		fmt.Printf("[INFO] [newAuthTestRouter] user created and added to the router: %s\n", auth.Id)
	}

	return testRouter{router: router, auth: auth, repos: &repos}
}

func (r testRouter) get(path string, body []byte) *httptest.ResponseRecorder {
	payload := bytes.NewReader(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, payload)

	if r.auth != nil {
		sessionId := uuid.New().String()
		// TODO: create redis test instance and add session to it
		req.AddCookie(&http.Cookie{
			Name:   "sessionId",
			Value:  sessionId,
			MaxAge: 3600 * 24,
		})

		err := r.repos.redis.Set(context.Background(), sessionId, r.auth.Id.String(), 10000).Err()

		if err != nil {
			panic(err)
		}

		fmt.Printf("[INFO] [testRouter.get] user id stored on redis: session(%s) => %s\n", sessionId, r.auth.Id)
	}

	r.router.ServeHTTP(w, req)

	return w
}
