package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type MeResponse struct {
	User *User `json:"user" tstype:"User"`
}

func (*API) user(ctx *gin.Context) (*User, error) {
	userJson := ctx.GetString("user")
	var user User
	err := json.Unmarshal([]byte(userJson), &user)

	return &user, err
}

func (s *API) handleMe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := s.user(ctx)

		if err != nil {
			fmt.Printf("[ERROR] [UserController.me] failed to get user from context: %s\n", err)
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		response := MeResponse{User: user}

		ctx.JSON(http.StatusOK, response)
	}
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required" tstype:"string"`
	Password string `json:"password" validate:"required,min=8,max=255" tstype:"string"`
}

type LoginResponse struct {
	User      *User  `json:"user" tstype:"User"`
	SessionId string `json:"sessionId" tstype:"string"`
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

		sessionId, err := login(user, ctx, s.repositories.redis)

		if err != nil {
			fmt.Printf("[ERROR] [UserController.login] failed to set redis session: %s\n", err)
			ctx.JSON(500, gin.H{"error": "unexpected error"})
			return
		}

		response := LoginResponse{User: user, SessionId: sessionId}
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
	User      *User  `json:"user" tstype:"User"`
	SessionId string `json:"sessionId" tstype:"string"`
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

		sessionId, err := login(user, ctx, s.repositories.redis)

		if err != nil {
			fmt.Printf("[ERROR] [UserController.login] failed to set redis session: %s\n", err)
			ctx.JSON(500, gin.H{"error": "unexpected error"})
			return
		}

		response := RegisterResponse{User: user, SessionId: sessionId}
		ctx.JSON(200, response)
	}
}

func (s *API) handleGoogleOAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientId, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID")
		if !ok {
			ctx.JSON(400, gin.H{"error": "Missing client id env"})
			return
		}

		redirectURL := "http://localhost:3000"
		authURL := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=https://www.googleapis.com/auth/userinfo.email", clientId, redirectURL)

		ctx.JSON(http.StatusOK, gin.H{"url": authURL})
	}
}

type GoogleOAuthExchangePayload struct {
	Code string `form:"code"`
}

func (s *API) handleGoogleOAuthExchange() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GoogleOAuthExchangePayload{}
		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Missing authorization code"})
			ctx.Abort()
			return
		}

		clientId, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_ID")
		if !ok {
			ctx.JSON(400, gin.H{"error": "Missing client id env"})
			return
		}

		clientSecret, ok := os.LookupEnv("GOOGLE_OAUTH_CLIENT_SECRET")
		if !ok {
			ctx.JSON(400, gin.H{"error": "Missing client id env"})
			return
		}

		redirectURL := "http://localhost:3000"
		tokenURL := "https://oauth2.googleapis.com/token"
		form := url.Values{}
		form.Set("client_id", clientId)
		form.Set("client_secret", clientSecret)
		form.Set("redirect_uri", redirectURL)
		form.Set("code", req.Code)
		form.Set("grant_type", "authorization_code")
		resp, err := http.PostForm(tokenURL, form)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch access token"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var tokenResponse map[string]interface{}
		err = json.Unmarshal(body, &tokenResponse)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse access token response body as JSON"})
			return
		}

		if tokenResponse["error"] != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": tokenResponse["error"].(string)})
			return
		}

		accessToken := tokenResponse["access_token"].(string)

		client := http.Client{}
		r, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v3/userinfo", nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create userinfo.email request"})
			return
		}
		r.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err = client.Do(r)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch userinfo.email"})
		}
		defer resp.Body.Close()

		var userData map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse userinfo response body as JSON"})
			return
		}

		userEmail := userData["email"].(string)
		user, err := s.repositories.userRepository.GetUserByEmail(userEmail)

		if user == nil || err != nil {
			user, err = s.repositories.userRepository.CreateOAuthUser(userEmail, strings.Split(userEmail, "@")[0])
			if err != nil {
				fmt.Printf("[ERROR] [UserController.handleGoogleOAuthExchange] failed to store user: %s\n", err)
				ctx.JSON(400, gin.H{"error": "Cannot create user"})
				return
			}
		}

		sessionId := uuid.New().String()

		err = s.repositories.redis.Set(ctx, sessionId, user.Id.String(), 0).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set session id"})
			return
		}

		ctx.SetCookie("sessionId", sessionId, 3600*24, "/", "localhost", false, false)

		ctx.JSON(http.StatusOK, gin.H{"user": user, "sessionId": sessionId})
	}
}

// TODO: move somewhere else
func login(user *User, ctx *gin.Context, redis *redis.Client) (string, error) {
	sessionId := uuid.New().String()

	err := redis.Set(ctx, sessionId, user.Id.String(), 0).Err()
	if err != nil {
		return "", err
	}

	ctx.SetCookie("sessionId", sessionId, 3600*24, "/", "localhost", false, false)

	return sessionId, err
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

		user, err := s.user(ctx)

		if err != nil {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

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
			ctx.Abort()
			return
		}

		user, err := s.user(ctx)

		if err != nil {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		if req.Limit == 0 {
			req.Limit = 10
		}

		resources, err := s.repositories.userRepository.GetUserResources(user, req.Limit, req.Offset, req.Status, req.ReviewRating)

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

		user, err := s.user(ctx)

		if err != nil {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		resource, err = s.repositories.userRepository.GetUserResource(user, resource.Id.String())

		if err != nil {
			ctx.JSON(200, gin.H{"resource": resource, "user_holds": false})
			return
		}

		ctx.JSON(200, gin.H{"resource": resource, "user_holds": true})
	}
}

func (s *API) handleGetRecommendations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := s.user(ctx)

		if err != nil {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		// TODO: get recommendations system URL from .env
		resp, err := http.Get("http://localhost:8000/" + user.Id.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recommendations from the system"})
			return
		}
		defer resp.Body.Close()

		var data interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response body as JSON"})
			return
		}

		ctx.JSON(http.StatusOK, data)
	}
}

func (s *API) handleGetPopularThisWeek() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results := new([]Resource)
		offset := 1

		for offset <= 4 {
			resources, err := s.repositories.resourceRepository.GetPopularThisWeekResources(offset)

			if err != nil {
				fmt.Printf("[ERROR] [API.GetPopularThisWeek] failed to fetch resources: %s\n", err)
				break
			}

			if len(*resources) <= 0 {
				offset += 1
				continue
			}

			results = resources
			break
		}

		fmt.Printf("[INFO] [API.GetPopularThisWeek] resources: %v\n", results)

		if len(*results) <= 0 {
			ctx.JSON(200, gin.H{"resources": []*Resource{}})
			ctx.Abort()
			return
		}

		ctx.JSON(200, gin.H{"resources": results})
	}
}
