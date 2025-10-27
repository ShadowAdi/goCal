package services

import (
	"errors"
	"goCal/internal/logger"
	"io"
)

type FileStorageService struct {
}

func NewFileStorageService() *FileStorageService {
	return &FileStorageService{}
}

func (nfs *FileStorageService) UploadSingleFile(userId string, fileName string, file io.Reader, fileType string) error {
	if userId == "" {
		logger.Error("Failed to get the userId UnAuthorized")
		return errors.New("Unauthorized User. UserId Not Found")
	}
	if fileName == "" {
		logger.Error("Failed to get the File Name")
		return errors.New("Failed to get the fileName")
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

	// TODO: Implement actual file upload logic using bucketName
	logger.Info("File will be uploaded to bucket: " + bucketName)

	return nil
}
