package routes

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func UserRoutes(router *gin.RouterGroup) {
	router.GET("/")
}
