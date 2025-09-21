package config

import (
	"fmt"
	"goCal/internal/logger"

	"github.com/joho/godotenv"
)

func GetLoadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		logger.Error(`Failed to load env vars ` + err.Error())
		fmt.Printf("Failed to load env vars: %s ", err.Error())
		return
	}
}
