package main

import (
	"log"

	auth "server/lib/middlewares/auth"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	db, err := sqlx.Connect("postgres", "user=postgres dbname=go-api password=password sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.Use(auth.AuthMiddleware(db, rdb))

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })

	router.Run("localhost:8080")
}
