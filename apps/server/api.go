package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
)

type API struct {
	addr         string
	repositories Repositories
}

func NewAPI(addr string, repositories Repositories) *API {
	return &API{
		addr:         addr,
		repositories: repositories,
	}
}

func (a *API) Run() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=banoffee password=5up3r_s3cur3_p4ssw0rd sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	router := gin.Default()
	router.Group("/").Use(AuthMiddleware(db, redis))

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })
	router.POST("/login", a.handleLogin())

	router.Run("localhost:8080")
}

func (s *API) handleLogin(ctx *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
