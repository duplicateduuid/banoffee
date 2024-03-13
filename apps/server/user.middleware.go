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
		session_id, err := ctx.Cookie("session_id")

		if err != nil {
			ctx.Next()
			return
		}

		user_id, err := rdb.Get(ctx, session_id).Result()

		if err != nil {
			ctx.Next()
			return
		}

		user_auth_info := UserAuthInfo{}
		database_error := db.Get(&user_auth_info, `SELECT u.id, u.email FROM "user" u WHERE u.id=$1`, user_id)
		if database_error != nil {
			ctx.Next()
			return
		}

		ctx.Set("user_id", user_auth_info.Id)
		ctx.Set("user_email", user_auth_info.Email)

		ctx.Next()
	}
}
