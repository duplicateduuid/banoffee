package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Login    string `json:"login" validate:"required" tstype:"string"`
	Password string `json:"password" validate:"required,min=8,max=255" tstype:"string"`
}

type LoginResponse struct {
	User *User `json:"user" tstype:"User"`
}

func (s *API) handleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := LoginRequest{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		validate := validator.New()

		err := validate.Struct(req)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			ctx.JSON(422, gin.H{"error": fmt.Sprintf("validation errors: %s", errors)})

			return
		}

		user, err := s.repositories.userRepository.GetUserByUsernameOrEmail(req.Login)
		if err != nil {
			ctx.JSON(403, gin.H{"error": "invalid email or password"})
			return
		}

		if user.ValidPassword(req.Password) {
			ctx.JSON(403, gin.H{"error": "invalid email or password"})
			return
		}

		session_id := uuid.New().String()
		err = s.repositories.redis.Set(ctx, session_id, user.Id.String(), 0).Err()
		if err != nil {
			fmt.Printf("[ERROR] [UserController.login] failed to set redis session: %s\n", err)
			ctx.JSON(500, gin.H{"error": "unexpected error"})
			return
		}

		ctx.SetCookie("sessionId", session_id, 3600*24, "/", "localhost", false, false)

		response := LoginResponse{User: user}
		ctx.JSON(200, response)
	}
}

type RegisterRequest struct {
	Email     string  `json:"email" validate:"required,email" tstype:"string"`
	Username  string  `json:"username" validate:"required,min=5,max=20" tstype:"string"`
	Password  string  `json:"password" validate:"required,min=8,max=255" tstype:"string"`
	AvatarUrl *string `json:"avatar_url" validate:"omitempty,http_url" tstype:"string | null"`
	HeaderUrl *string `json:"header_url" validate:"omitempty,http_url" tstype:"string | null"`
	Bio       *string `json:"bio" validate:"omitempty,max=255" tstype:"string | null"`
}

type RegisterResponse struct {
	User *User `json:"user" tstype:"User"`
}

func (s *API) hanlderRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := RegisterRequest{}

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

		user, err := NewUser(req.Email, req.Username, req.Password, req.AvatarUrl, req.HeaderUrl, req.Bio)

		if err != nil {
			fmt.Printf("[ERROR] [UserController.register] failed to create user: %s", err)
			ctx.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err = s.repositories.userRepository.CreateUser(user)
		if err != nil {
			fmt.Printf("[ERROR] [UserController.register] failed to store user: %s\n", err)
			ctx.JSON(400, gin.H{"error": "Cannot create user"})
			return
		}

		response := RegisterResponse{User: user}
		ctx.JSON(200, response)
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

		// TODO: this is bad. Other errors besides the duplicated one can happen.
		if err != nil {
			err = s.repositories.userRepository.CreateUserResource(user, resourceId, req.Status, req.ReviewRating, req.ReviewComment)

			if err != nil {
				fmt.Printf("[ERROR] [UserController.saveResource] failed to create user resource: %s\n", err)
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
				fmt.Printf("[ERROR] [API.SaveResource] failed to update resource: %s\n", err)
				ctx.JSON(400, gin.H{"error": "Failed to update resource"})
				return
			}

			ctx.JSON(200, gin.H{"message": "Resource updated with success"})
		}
	}
}

type GetMyResourcesPayload struct {
	Limit  int `db:"limit" form:"limit"`
	Offset int `db:"offset" form:"offset"`
}

func (s *API) handleGetMyResources() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetMyResourcesPayload{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			ctx.Abort()
			return
		}

		user := ctx.MustGet("user").(*User)

		if req.Limit == 0 {
			req.Limit = 10
		}

		resources, err := s.repositories.resourceRepository.GetUserResources(user, req.Limit, req.Offset)

		if err != nil {
			fmt.Printf("[ERROR] [API.GetMyResources] failed to fetch resources: %s", err)
			ctx.JSON(400, gin.H{"error": "Cannot retrieve resources"})
			ctx.Abort()
			return
		}

		if len(*resources) <= 0 {
			ctx.JSON(200, gin.H{"resources": []*Resource{}})
			ctx.Abort()
			return
		}

		ctx.JSON(200, gin.H{"resources": resources})
	}
}

type GetMyResourcePayload struct {
	Url string `db:"url" form:"url"`
}

// TODO: handle possible errors here instead of just return null for all error cases
func (s *API) handleGetMyResource() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetMyResourcePayload{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid URL"})
			ctx.Abort()
			return
		}

		resource, err := s.repositories.resourceRepository.GetResourceByUrl(req.Url)

		if err != nil {
			ctx.JSON(200, gin.H{"resource": nil})
			return
		}

		user := ctx.MustGet("user").(*User)

		_, err = s.repositories.userRepository.GetUserResource(user, resource.Id.String())

		if err != nil {
			ctx.JSON(200, gin.H{"resource": resource, "user_holds": false})
			return
		}

		ctx.JSON(200, gin.H{"resource": resource, "user_holds": true})
	}
}
