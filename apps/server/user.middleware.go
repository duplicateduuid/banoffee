package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type UserAuthInfo struct {
	Id    uuid.UUID `db:"id" json:"id"`
	Email string    `db:"email" json:"email"`
}

func AuthMiddleware(db *sqlx.DB, rdb *redis.Client) gin.HandlerFunc {
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

		userAuthInfo := UserAuthInfo{}
		err = db.Get(&userAuthInfo, `SELECT u.id, u.email FROM "user" u WHERE u.id=$1`, userId)
		if err != nil {
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("userId", userAuthInfo.Id)
		ctx.Set("userEmail", userAuthInfo.Email)

		ctx.Next()
	}
}
