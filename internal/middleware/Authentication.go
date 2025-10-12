package middleware

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/types"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	JWT_KEY := os.Getenv("JWT_KEY")
	if JWT_KEY == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("token")
		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "Missing Authorization header", "success": false})
			ctx.Abort()
			return
		}

		claims := &types.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return JWT_KEY, nil
			})
		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"error": "Invalid token", "success": false})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.ID)
		ctx.Set("email", claims.Issuer)
		ctx.Next()
	}
}
