package main

import (
	"aegis-ai-gateway/internal/config"
	"aegis-ai-gateway/internal/logging"
	"aegis-ai-gateway/internal/server"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		panic("Failed to load config " + err.Error())
	}
	appLogger := logging.InitLogger()
	appContext := context.Background()
	appLogger.Log(appContext, slog.LevelInfo, "Starting Aegis AI Gateway")
	appLogger.Info("Config is :" + config.App_Env + config.Database_URL + config.Log_Format + config.Log_Level + config.Request_Timeout)
	appServer := server.NewServer(server.Config{Addr: config.Http_Addr}, appLogger)
	// convert the the app server into a real net/http server
	httpServer := appServer.HttpServer()

	// create a context that is canceleld when the process receives siginti
	ctx, _ := signal.NotifyContext(appContext, syscall.SIGINT, syscall.SIGTERM)

	// Start the http server in a goroutine
	//
	go func() {
		appLogger.Info(
			"Starintng api",
			"Service", "aegis-ai-gateway",
			"addr", config.Http_Addr,
			"env", config.App_Env,
		)

		if err := httpServer.ListenAndServe(); err != nil {
			appLogger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()

	appLogger.Info("Shutdown signal recieved")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Failed to shutdown server", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Api Server Stopped")

}
