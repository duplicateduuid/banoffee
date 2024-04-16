package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type API struct {
	repositories Repositories
}

func NewAPI(reops Repositories) *API {
	return &API{
		repositories: reops,
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: use .env for origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (a *API) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(CorsMiddleware())

	authRouter := router.Group("/").Use(a.AuthMiddleware())

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })
	router.POST("/login", a.handleLogin())
	router.POST("/register", a.hanlderRegister())
	router.GET("/popular", a.handleGetPopularThisWeek())

	authRouter.GET("/me", a.handleMe())
	authRouter.GET("/resource", a.handleGetResource())
	authRouter.POST("/resource", a.handleCreateResource())
	authRouter.GET("/resource/search", a.handleSearchResource())

	authRouter.GET("/user/resources", a.handleGetMyResources())
	authRouter.GET("/user/resource", a.handleGetMyResource())
	authRouter.POST("/user/resource/:id", a.handleSaveResource())

	authRouter.GET("/recommendations", a.handleGetRecommendations())

	return router
}

func (a *API) Run(addr string) {
	router := a.SetupRouter()

	router.Run(addr)
}
