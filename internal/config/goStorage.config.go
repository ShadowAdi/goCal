package config

import (
	"fmt"
	"goCal/internal/logger"
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

var storageClient *storage_go.Client

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

}
