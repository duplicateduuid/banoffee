package main

import (
	"fmt"

	"encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type API struct {
	addr         string
	repositories Repositories
	redis        *redis.Client
}

func NewAPI(addr string, repositories Repositories) *API {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &API{
		addr:         addr,
		repositories: repositories,
		redis:        redis,
	}
}

func (a *API) Run() {
	router := gin.Default()
	router.Group("/").Use(AuthMiddleware(a.repositories.userRepository, a.redis))

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })
	router.POST("/login", a.handleLogin())
	router.POST("/register", a.hanlderRegister())

	router.Run(a.addr)
}

type loginJson struct {
	email    string
	password string
}

func (s *API) handleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload loginJson

		if ctx.ShouldBindJSON(&payload) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		var user User
		// TODO: validate email
		if s.repositories.userRepository.GetUserByEmail(payload.email, user) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid email or password"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.password)) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid email or password"})
			return
		}

		body, err := json.Marshal(user)

		if err != nil {
			ctx.JSON(400, gin.H{"error": "failed to serialize user"})
		}

		ctx.JSON(200, gin.H{"user": body})
	}
}

func (s *API) hanlderRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := User{}

		if ctx.ShouldBindJSON(&payload) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		err := s.repositories.userRepository.CreateUser(&payload)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Cannot create user"})
			return
		}

		ctx.JSON(200, gin.H{"message": "User created with success!"})
	}
}
