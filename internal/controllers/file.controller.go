package controllers

import "github.com/gin-gonic/gin"

type FileController struct{}

func NewFileController() *UserController {
	return &UserController{}
}

func (fc *FileController) GetAllFiles(ctx *gin.Context) {

}
