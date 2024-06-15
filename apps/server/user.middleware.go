package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (a *API) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionId, err := ctx.Cookie("sessionId")

		if err != nil {
			fmt.Printf("[ERROR] [AuthMiddleware] failed to get session id from cookies: %s\n", err)
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		valkey := a.repositories.valkey
		userId, err := valkey.Do(context.Background(), valkey.B().Get().Key(sessionId).Build()).ToString()

		if err != nil {
			fmt.Printf("[ERROR] [AuthMiddleware] session(%s) not found on valkey: %s\n", sessionId, err)
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		id, err := uuid.Parse(userId)

		if err != nil {
			fmt.Printf("[ERROR] [AuthMiddleware] failed to parse user id: %s\n", err)
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		user, err := a.repositories.userRepository.GetUserById(id)
		if err != nil {
			fmt.Printf("[ERROR] [AuthMiddleware] failed to get user by id: %s\n", err)
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}
		b, err := json.Marshal(user)

		if err != nil {
			fmt.Printf("[ERROR] [AuthMiddleware] failed to fetch user: %s\n", err)
			ctx.JSON(401, gin.H{"message": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Set("user", string(b))

		ctx.Next()
	}
}
