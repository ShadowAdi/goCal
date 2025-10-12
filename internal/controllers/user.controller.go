package controllers

import (
	"goCal/internal/schema"
	"goCal/internal/services"
	"goCal/internal/utils"
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
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
	return
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	user, error := uc.UserService.GetUser(id)
	if error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var newUser *schema.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err,
		})
		return
	}
	newUser.Password = hashedPassword

	user, error := uc.UserService.CreateUser(newUser)
	if error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   error,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   user,
	})
	return
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	message, err := uc.UserService.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
	})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	_, userFoundError := uc.UserService.GetUser(id)
	if userFoundError != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   userFoundError.Error(),
		})
		return
	}

	var updateRequest *schema.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&updateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updatedUser, updateUserError := uc.UserService.UpdateUser(id, updateRequest)
	if updateUserError != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   updateUserError.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    updatedUser,
	})
}
