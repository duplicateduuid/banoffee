package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type API struct {
	repositories Repositories
	redis        *redis.Client
}

func NewAPI(repositories Repositories) *API {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &API{
		repositories: repositories,
		redis:        redis,
	}
}

func (a *API) SetupRouter() *gin.Engine {
	router := gin.Default()
	authRouter := router.Group("/").Use(AuthMiddleware(a.repositories.userRepository, a.redis))

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })
	router.POST("/login", a.handleLogin())
	router.POST("/register", a.hanlderRegister())

	authRouter.GET("/resource", a.handleGetResource())
	authRouter.POST("/resource", a.handleCreateResource())
	authRouter.GET("/resources", a.handleGetMyResources())

	return router
}

func (a *API) Run(addr string) {
	router := a.SetupRouter()

	router.Run(addr)
}
