package controllers

import (
	"fmt"
	"goCal/internal/logger"
	"goCal/internal/schema"
	"goCal/internal/services"
	"goCal/internal/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserController struct {
	UserService *services.UserService
}

var ADMIN_EMAIL string

func Init() {
	ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	if ADMIN_EMAIL == "" {
		fmt.Printf(`Failed to get the ADMIN_EMAIL`)
	}
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
		"message": "User created successfully. Please check your email for verification code.",
		"user":    user,
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

	// Check if user is verified
	if !userFound.IsVerified {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Please verify your email before logging in",
			"user_id": userFound.ID.String(),
		})
		return
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userFound.Email,
		Id:        userFound.ID.String(),
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
	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User Id not found in context",
		})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	_, loggedInUserError := uc.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		fmt.Printf("Error finding logged-in user: %v\n", loggedInUserError)
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	message, err := uc.UserService.DeleteUser(userIdStr)
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
	email, exists := ctx.Get("email")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"success": false,
			"error":   "Email not found in context",
		})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid email type in context",
		})
		return
	}

	loggedInUserFound, loggedInUserError := uc.UserService.GetUserByEmail(emailStr)
	if loggedInUserError != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
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

	updatedUser, updateUserError := uc.UserService.UpdateUser(loggedInUserFound.ID.String(), updateRequest)
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

	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User Id not found in context",
		})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	loggedInUser, loggedInUserError := uc.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	if strings.ToLower(loggedInUser.Email) != ADMIN_EMAIL {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Only Admin Can Access this api",
		})
		return
	}

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

	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User Id not found in context",
		})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	loggedInUser, loggedInUserError := uc.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	if strings.ToLower(loggedInUser.Email) != ADMIN_EMAIL {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Only Admin Can Access this api",
		})
		return
	}

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

	userId, exists := ctx.Get("userId")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User Id not found in context",
		})
		return
	}

	userIdStr, ok := userId.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid User Id type in context",
		})
		return
	}

	loggedInUser, loggedInUserError := uc.UserService.GetUser(userIdStr)
	if loggedInUserError != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   loggedInUserError.Error(),
		})
		return
	}

	if strings.ToLower(loggedInUser.Email) != ADMIN_EMAIL {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Only Admin Can Access this api",
		})
		return
	}

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

// ResendVerificationEmail resends verification email to a user
func (uc *UserController) ResendVerificationEmail(ctx *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	emailResponse, err := uc.UserService.ResendVerificationEmail(request.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": emailResponse.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": emailResponse.Message,
	})
}

// VerifyUser verifies a user with the provided verification code
func (uc *UserController) VerifyUser(ctx *gin.Context) {
	var request struct {
		Email            string `json:"email" binding:"required,email"`
		VerificationCode string `json:"verification_code" binding:"required,len=4"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	user, err := uc.UserService.VerifyUser(request.Email, request.VerificationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User verified successfully",
		"user":    user,
	})
}
