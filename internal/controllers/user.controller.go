package controllers

import (
	"goCal/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	users, error := uc.UserService.GetUsers()
	if error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   error,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
	return
}

func (uc *UserController) GetUser(id string, ctx *gin.Context) {
	user, error := uc.UserService.GetUser(id)
	if error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   error,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   user,
	})
	return
}
