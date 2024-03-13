package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// _, err := sqlx.Connect("postgres", "user=postgres dbname=go-api password=password sslmode=disable")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	router := gin.Default()

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Banoffee"}) })

	router.Run("localhost:8080")
}
