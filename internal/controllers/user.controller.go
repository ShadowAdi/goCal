package controllers

import (
	"fmt"
	"goCal/internal/db"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	var users []schema.User
	result := db.DB.Find(&users)

	if result.Error != nil {
		logger.Error(`Failed to get Users %w`, result.Error)
		panic(fmt.Errorf("Failed to get all users: %w", result.Error))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
	return
}

func (uc *UserController) GetUser(id string, ctx *gin.Context) {
	var user schema.User
	result := db.DB.Where("id = ?", id).First(&user)

	if result.Error != nil {
		logger.Error(`Failed to get User %w`, result.Error)
		panic(fmt.Errorf("Failed to get user: %w", result.Error))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   user,
	})
	return
}
