package middleware

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/types"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	var ADMIN_EMAIL string

	ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	if ADMIN_EMAIL == "" {
		fmt.Printf(`Failed to get the ADMIN_EMAIL`)
	}
	JWT_KEY := os.Getenv("JWT_KEY")
	if JWT_KEY == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}
	return func(ctx *gin.Context) {
		var tokenString string

		// Check Authorization header first (standard way)
		authHeader := ctx.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>" format
			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
				tokenString = tokenParts[1]
			}
		}

		// Fallback to token header if Authorization is not present or invalid
		if tokenString == "" {
			tokenString = ctx.GetHeader("token")
		}

		if tokenString == "" {
			ctx.JSON(401, gin.H{"error": "Missing Authorization header", "success": false})
			ctx.Abort()
			return
		}

		claims := &types.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims,
			func(t *jwt.Token) (interface{}, error) {
				return []byte(JWT_KEY), nil
			})
		if err != nil {
			fmt.Printf("JWT Parse Error: %v", err)
			ctx.JSON(401, gin.H{"error": "Invalid token", "success": false})
			ctx.Abort()
			return
		}
		if !token.Valid {
			fmt.Printf("Token is not valid")
			ctx.JSON(401, gin.H{"error": "Invalid token", "success": false})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Id)
		ctx.Set("email", claims.Issuer)
		if strings.ToLower(claims.Issuer) == strings.ToLower(ADMIN_EMAIL) {
			ctx.Set("role", "admin")
		} else {
			ctx.Set("role", "user")
		}

		ctx.Next()
	}
}
