package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginPayload struct {
	Email    string `db:"email" json:"email"`
	Password string `json:"password"`
}

func (s *API) handleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := LoginPayload{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		// TODO: validate email
		user, err := s.repositories.userRepository.GetUserByEmail(req.Email)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid email or password"})
			return
		}

		if user.ValidPassword(req.Password) {
			ctx.JSON(400, gin.H{"error": "Invalid email or password"})
			return
		}

		body, err := json.Marshal(user)

		if err != nil {
			ctx.JSON(400, gin.H{"error": "Failed to serialize user"})
		}

		session_id := uuid.New().String()
		err = s.redis.Set(ctx, session_id, user.Id.String(), 0).Err()
		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Unexpected error"})
			return
		}

		ctx.JSON(200, gin.H{"user": body})
	}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
	HeaderUrl string `json:"header_url"`
	Bio       string `json:"bio"`
}

func (s *API) hanlderRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := RegisterRequest{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		user, err := NewUser(req.Email, req.Username, req.Password, req.AvatarUrl, req.HeaderUrl, req.Bio)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		if s.repositories.userRepository.CreateUser(user) != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Cannot create user"})
			return
		}

		ctx.JSON(200, gin.H{"message": "User created with success!"})
	}
}