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

func ensureBucket(name string) {
	_, err := storageClient.GetBucket(name)
	if err == nil {
		fmt.Printf("Bucket %s already exists\n", name)
		return
	}

	// If not found, create it
	_, err = storageClient.CreateBucket(name, storage_go.BucketOptions{
		Public:        true,
		FileSizeLimit: "1000",
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create bucket %s: %v", name, err))
		return
	}
	fmt.Printf("Bucket %s created successfully\n", name)
}

func StorageInit() {
	url := os.Getenv("SUPABASE_PROJECT_URL")
	key := os.Getenv("SUPABASE_PROJECT_KEY")

	if url == "" || key == "" {
		logger.Error("Missing Supabase credentials")
		return
	}

	storageClient = storage_go.NewClient(url, key, nil)

	ensureBucket("goCal-Other-Bucket")
	ensureBucket("goCal-Docs-Bucket")
	ensureBucket("goCal-Albums-Bucket")
	ensureBucket("goCal-Videos-Bucket")
}
