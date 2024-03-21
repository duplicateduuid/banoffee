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

func (a *API) SetupRouter() *gin.Engine {
	router := gin.Default()
	authRouter := router.Group("/").Use(a.AuthMiddleware())

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })
	router.POST("/login", a.handleLogin())
	router.POST("/register", a.hanlderRegister())

	authRouter.GET("/resource", a.handleGetResource())
	authRouter.POST("/resource", a.handleCreateResource())

	authRouter.GET("/resource/search", a.handleSearchResource())
	authRouter.GET("/user/resource", a.handleGetMyResources())
	authRouter.GET("/user/resource/:id", a.handleGetMyResource())
	authRouter.POST("/user/resource/:id", a.handleSaveResource())

	return router
}

func (a *API) Run(addr string) {
	router := a.SetupRouter()

	router.Run(addr)
}
