package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoginPayload struct {
	Email    string `db:"email" json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (s *API) handleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := LoginPayload{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		validate := validator.New()

		err := validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("validation errors: %s", errors)})
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

		session_id := uuid.New().String()
		err = s.redis.Set(ctx, session_id, user.Id.String(), 0).Err()
		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Unexpected error"})
			return
		}

		ctx.SetCookie("sessionId", session_id, 3600*24, "/", "localhost", false, true)
		ctx.JSON(200, gin.H{"user": user})
	}
}

type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email"`
	Username  string  `json:"username" validate:"required,min=5,max=20"`
	Password  string  `json:"password" validate:"required,min=8,max=255"`
	AvatarUrl *string `json:"avatar_url" validate:"omitempty,http_url"`
	HeaderUrl *string `json:"header_url" validate:"omitempty,http_url"`
	Bio       *string `json:"bio" validate:"omitempty,max=255"`
}

func (s *API) hanlderRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := RegisterRequest{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		fmt.Printf("%v", req)

		validate := validator.New()
		err := validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			ctx.JSON(400, gin.H{"error": fmt.Sprintf("validation errors: %s", errors)})
			return
		}

		user, err := NewUser(req.Email, req.Username, req.Password, req.AvatarUrl, req.HeaderUrl, req.Bio)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		err = s.repositories.userRepository.CreateUser(user)
		if err != nil {
			fmt.Printf("[ERROR]: %v\n", err)
			ctx.JSON(400, gin.H{"error": "Cannot create user"})
			return
		}

		ctx.JSON(200, gin.H{"message": "User created with success!"})
	}
}

type GetMyResourcesPayload struct {
	Limit        int    `db:"limit" form:"limit"`
	Offset       int    `db:"offset" form:"offset"`
	Status       string `db:"status" form:"status"`
	ReviewRating string `db:"review_rating" form:"review_rating"`
}

func (s *API) handleGetMyResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetMyResourcesPayload{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		if req.Limit == 0 {
			req.Limit = 10
		}

		user := ctx.MustGet("user").(*User)

		resources, err := s.repositories.userRepository.GetUserResources(user, req.Limit, req.Offset, req.Status, req.ReviewRating)

		if err != nil {
			fmt.Println(err)
			ctx.JSON(400, gin.H{"error": "Cannot retrieve resources"})
			return
		}

		if len(*resources) <= 0 {
			ctx.JSON(200, gin.H{"resources": []*Resource{}})
			return
		}

		ctx.JSON(200, gin.H{"resources": resources})
	}
}

type SaveResourcePayload struct {
	Status        *string `db:"status" json:"status"`
	ReviewRating  *string `db:"review_rating" json:"review_rating"`
	ReviewComment *string `db:"review_comment" json:"review_comment"`
}

func (s *API) handleSaveResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resourceId := ctx.Param("id")

		if resourceId == "" {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		req := SaveResourcePayload{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		user := ctx.MustGet("user").(*User)

		resource, err := s.repositories.resourceRepository.GetResourceById(resourceId)

		if err != nil {
			ctx.JSON(400, gin.H{"error": "Cannot retrieve resource"})
			return
		}

		_, err = s.repositories.userRepository.GetUserResource(user, resourceId)

		if err != nil {
			err = s.repositories.userRepository.CreateUserResource(user, resourceId, req.Status, req.ReviewRating, req.ReviewComment)

			if err != nil {
				fmt.Println(err)
				ctx.JSON(400, gin.H{"error": "Failed to create resource"})
				return
			}

			ctx.JSON(200, gin.H{"message": "Resource created with success"})
		} else {
			newStatus := req.Status

			if resource.Status == nil && req.Status == nil && req.ReviewRating != nil {
				updatedStatus := "ongoing"
				newStatus = &updatedStatus
			}

			err = s.repositories.userRepository.UpdateUserResource(user, resourceId, newStatus, req.ReviewRating, req.ReviewComment)

			if err != nil {
				fmt.Println(err)
				ctx.JSON(400, gin.H{"error": "Failed to update resource"})
				return
			}

			ctx.JSON(200, gin.H{"message": "Resource updated with success"})
		}
	}
}
