package services

import (
	"bytes"
	"errors"
	"fmt"
	"goCal/internal/logger"
	"io"
	"path/filepath"
	"time"

	storage_go "github.com/supabase-community/storage-go"
)

type FileStorageService struct {
	storageClient *storage_go.Client
}

func NewFileStorageService(client *storage_go.Client) *FileStorageService {
	return &FileStorageService{
		storageClient: client,
	}
}

func (nfs *FileStorageService) UploadFile(userId string, fileName string, file io.Reader, fileType string) (string, error) {
	if userId == "" {
		logger.Error("Failed to get the userId UnAuthorized")
		return "", errors.New("Unauthorized User. UserId Not Found")
	}
	if fileName == "" {
		logger.Error("Failed to get the File Name")
		return "", errors.New("Failed to get the fileName")
	}

	var bucketName string
	switch fileType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/bmp", "image/webp":
		bucketName = "goCal-Albums-Bucket"
	case "video/mp4", "video/webm", "video/avi", "video/mov", "video/mkv":
		bucketName = "goCal-Videos-Bucket"
	case "audio/mp3", "audio/wav", "audio/aac", "audio/ogg", "audio/flac":
		bucketName = "goCal-Audios-Bucket"
	case "application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		bucketName = "goCal-Docs-Bucket"
	default:
		bucketName = "goCal-Other-Bucket"
	}

	timeStamp := time.Now().Unix()
	fileExt := filepath.Ext(fileName)
	baseFileName := fileName[:len(fileName)-len(fileExt)]
	uniqueFileName := fmt.Sprintf("%s%d_%s%s", userId, timeStamp, baseFileName, fileExt)
	logger.Info(fmt.Sprintf("Uploading file %s to bucket: %s", uniqueFileName, bucketName))

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read file: %v", err))
		return "", errors.New("failed to read file: %w " + err.Error())
	}

	_, errUpload := nfs.storageClient.UploadFile(bucketName, uniqueFileName, bytes.NewReader(fileBytes))
	if errUpload != nil {
		logger.Error(fmt.Sprintf("Failed to upload file to supabase %v ", err))
		return fmt.Sprintf("Failed to upload file to supabase %v ", err), nil
	}

	publicURL := nfs.storageClient.GetPublicUrl(bucketName, uniqueFileName)
	logger.Info(fmt.Sprintf("File uploaded successfully. URL: %s", publicURL.SignedURL))

	return publicURL.SignedURL, nil
}
