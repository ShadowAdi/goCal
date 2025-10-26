package config

import (
	"fmt"
	"goCal/internal/logger"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

var storageClient *storage_go.Client
var storageBucketVideos storage_go.Bucket
var storageBucketAlbums storage_go.Bucket
var storageBucketDocs storage_go.Bucket
var storageBucketOthers storage_go.Bucket

func StorageInit() {
	SUPABASE_PROJECT_URL := os.Getenv("SUPABASE_PROJECT_URL")
	if SUPABASE_PROJECT_URL == "" {
		logger.Error(`Failed to get the SUPABASE Project Url`)
		fmt.Printf(`Failed to get the SUPABASE Project Url`)
	}

	SUPABASE_PROJECT_KEY := os.Getenv("SUPABASE_PROJECT_KEY")
	if SUPABASE_PROJECT_KEY == "" {
		logger.Error(`Failed to get the SUPABASE Project Key`)
		fmt.Printf(`Failed to get the SUPABASE Project Key`)
	}

	storageClient = storage_go.NewClient(SUPABASE_PROJECT_URL, SUPABASE_PROJECT_KEY, nil)

	var storageBucketAlbumError error
	storageBucketAlbums, storageBucketAlbumError = storageClient.CreateBucket("goCal-Albums-Bucket", storage_go.BucketOptions{
		Public:        true,
		FileSizeLimit: "10",
	})
	if storageBucketAlbumError != nil {
		logger.Error("Failed to create Albums bucket: " + storageBucketAlbumError.Error())
		fmt.Printf("Failed to create Albums bucket: %v", storageBucketAlbumError)
	}

	var storageBucketDocsError error
	storageBucketDocs, storageBucketDocsError = storageClient.CreateBucket("goCal-Docs-Bucket", storage_go.BucketOptions{
		Public:        true,
		FileSizeLimit: "10",
	})
	if storageBucketDocsError != nil {
		logger.Error("Failed to create Albums bucket: " + storageBucketDocsError.Error())
		fmt.Printf("Failed to create Albums bucket: %v", storageBucketDocsError)
	}

}
