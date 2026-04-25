package main

import (
	"aegis-ai-gateway/internal/config"
	"aegis-ai-gateway/internal/logging"
	"context"
	"log/slog"
)



func main() {
    config, err := config.GetConfig()
	if err != nil {
		panic("Failed to load config " + err.Error())
	}
	appLogger := logging.InitLogger()
	context := context.Background()
	appLogger.Log(context, slog.LevelInfo,"Starting Aegis AI Gateway",)
	appLogger.Info("Config is :" + config.App_Env + config.Database_URL + config.Log_Format + config.Log_Level + config.Request_Timeout)
	
	
	
	
}