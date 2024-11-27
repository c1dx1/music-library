package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"music-library/config"
	_ "music-library/docs"
	"music-library/internal/db"
	"music-library/internal/handlers"
	"music-library/internal/repositories"
	"music-library/internal/services"
	"music-library/internal/utils"
)

// @title Music library
// @version 1.0

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	logger := utils.LoadLogger(cfg)
	logger.Info("Logger initialized successfully")

	if err = db.RunMigrations(cfg.PostgresURL(), logger); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}
	logger.Info("Database migrations run successfully")

	pool, err := db.LoadDB(cfg.PostgresURL(), logger)
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
	}
	logger.Info("Connected to the database successfully")
	defer func() {
		logger.Debug("Closing database connection")
		pool.Close()
	}()

	logger.Debug("Initializing external API client")
	externalAPIClient := services.NewExternalAPIClient(cfg, logger)
	logger.Debug("External API client initialized successfully")

	logger.Debug("Initializing music library repository")
	mlibRepo := repositories.NewMLibRepository(pool, logger)
	logger.Debug("Music library repository initialized successfully")

	logger.Debug("Initializing music library service")
	mlibService := services.NewMLibService(*mlibRepo, logger, externalAPIClient)
	logger.Debug("Music library service initialized successfully")

	logger.Debug("Initializing handlers")
	handler := handlers.NewMLibHandler(mlibService, logger)
	logger.Debug("Handlers initialized successfully")

	logger.Debug("Setting up router")
	router := gin.Default()

	logger.Debug("Defining routes")

	router.GET("/songs", handler.GetLibrary)
	router.GET("/songs/:id", handler.GetText)
	router.DELETE("/songs/:id", handler.DeleteSong)
	router.PUT("/songs/:id", handler.EditSong)
	router.POST("/songs", handler.AddSong)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("API is listening on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerAddress + cfg.ServerPort); err != nil {
		logger.Fatalf("Could not start API: %v", err)
	}
}
