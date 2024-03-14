package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(repo UserRepository, rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := ctx.Cookie("sessionId")

		if err != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		userId, err := rdb.Get(ctx, sessionId).Result()

		if err != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		id, err := uuid.Parse(userId)

		if err != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		user, err := repo.GetUserById(id)

		if err != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	}
}
