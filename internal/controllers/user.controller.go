package controllers

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"goCal/internal/services"
	"goCal/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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

func (uc *UserController) LoginUser(ctx *gin.Context) {

	var newUser *schema.User
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	JWT_KEY := os.Getenv("JWT_KEY")
	if JWT_KEY == "" {
		logger.Error(`Failed to get the database url`)
		fmt.Printf(`Failed to get the database url`)
	}
	userFound, error := uc.UserService.GetUserByEmail(newUser.Email)
	if error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   error.Error(),
		})
		return
	}
	err := utils.CompareHashAndPassword(userFound.Password, newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    newUser.Email,
		Id:        newUser.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(JWT_KEY))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
		"email":   userFound.Email,
		"id":      userFound.ID,
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

// GetSoftDeletedUsers returns all soft-deleted users
func (uc *UserController) GetSoftDeletedUsers(ctx *gin.Context) {
	users, err := uc.UserService.GetSoftDeletedUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
}

// RestoreUser restores a soft-deleted user
func (uc *UserController) RestoreUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	user, err := uc.UserService.RestoreUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User restored successfully",
		"user":    user,
	})
}

// PermanentlyDeleteUser permanently deletes a user (hard delete)
func (uc *UserController) PermanentlyDeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "User ID is required",
		})
		return
	}

	err := uc.UserService.PermanentlyDeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User permanently deleted",
	})
}
