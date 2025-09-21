package main

import (
	"goCal/internal/config"
	"goCal/internal/logger"
)

func main() {
	logger.InitLogger()
	config.GetLoadEnvVars()
}
