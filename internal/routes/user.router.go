package routes

import (
	"fmt"
	"goCal/internal/db"
	"goCal/internal/schema"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/", func(ctx *gin.Context) {
		var user *schema.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err, "success": false})
			return
		}

		var isExists bool

		checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`

		err := db.Conn.QueryRow(ctx, checkQuery, user.Email).Scan(&isExists)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "success": false})
			return
		}
		if isExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered", "success": false})
			return
		}
		custom_link := strings.Split(user.Username, " ")[0] + "-" + strings.Split(user.Email, " ")[0]

		query := `
		INSERT INTO users (username, email, password, profileUrl, country, pronouns,custom_link)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id,custom_link,created_at,welcome_message,isVerified,date_format,time_format,timezone;
		`
		hashedPassword, err := hashedPassword(user.Password)
		if err != nil {
			fmt.Println("Error hashing password:", err)
			ctx.JSON(http.StatusUnauthorized,
				gin.H{"error": "Error hashing password", "success": false})
			return
		}

		row := db.Conn.QueryRow(ctx, query, user.Username, user.Email, hashedPassword, user.ProfileUrl, user.Country, user.Pronouns, custom_link)

		if err := row.Scan(&user.ID, &user.CustomLink, &user.CreatedAt, &user.WelcomeMessage, &user.IsVerified, &user.DateFormat, &user.TimeFormat, &user.Timezone); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"user": gin.H{
				"id":              user.ID,
				"username":        user.Username,
				"email":           user.Email,
				"country":         user.Country,
				"welcome_message": user.WelcomeMessage,
				"timezone":        user.Timezone,
				"pronouns":        user.Pronouns,
				"isverified":      user.IsVerified,
				"date_format":     user.DateFormat,
				"time_format":     user.TimeFormat,
				"custom_link":     user.CustomLink,
				"created_at":      user.CreatedAt,
			},
			"success": true,
		})

	})
	router.GET("/", func(ctx *gin.Context) {
		var users []schema.User

		query := `SELECT id,username,email,profileUrl,country,welcome_message,timezone,pronouns,date_format,time_format,custom_link,created_at FROM users`

		rows, err := db.Conn.Query(ctx, query)
		if err != nil {
			fmt.Println("Failed to get all users:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all users", "success": false})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user schema.User
			err := rows.Scan(
				&user.ID,
				&user.Username,
				&user.Email,
				&user.ProfileUrl,
				&user.Country,
				&user.WelcomeMessage,
				&user.Timezone,
				&user.Pronouns,
				&user.DateFormat,
				&user.TimeFormat,
				&user.CustomLink,
				&user.CreatedAt,
			)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users":   users,
			"success": true,
		})

	})
}
