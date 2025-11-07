package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"flypro-assessment-ayo/internal/config"
	"flypro-assessment-ayo/internal/handlers"
	"flypro-assessment-ayo/internal/middleware"
	"flypro-assessment-ayo/internal/repository"
	"flypro-assessment-ayo/internal/services"
	"flypro-assessment-ayo/internal/validators"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.App.Environment == "production" {
		logger, _ = zap.NewProduction()
		gin.SetMode(gin.ReleaseMode)
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Connect to database
	db, err := connectDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize validator
	validate := validator.New()
	if err := validators.RegisterCustomValidators(validate); err != nil {
		logger.Fatal("Failed to register validators", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	reportRepo := repository.NewExpenseReportRepository(db)

	// Initialize services
	currencyService := services.NewCurrencyService(redisClient, cfg)
	userService := services.NewUserService(userRepo, redisClient, cfg)
	expenseService := services.NewExpenseService(expenseRepo, currencyService, redisClient, cfg)
	reportService := services.NewExpenseReportService(reportRepo, expenseRepo, currencyService, redisClient, cfg)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, validate)
	expenseHandler := handlers.NewExpenseHandler(expenseService, validate)
	reportHandler := handlers.NewExpenseReportHandler(reportService, validate)

	// Setup router
	router := setupRouter(logger, userHandler, expenseHandler, reportHandler)

	// Setup server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func setupRouter(
	logger *zap.Logger,
	userHandler *handlers.UserHandler,
	expenseHandler *handlers.ExpenseHandler,
	reportHandler *handlers.ExpenseReportHandler,
) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler(logger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// User routes
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
		}

		// Expense routes
		expenses := api.Group("/expenses")
		{
			expenses.POST("", expenseHandler.CreateExpense)
			expenses.GET("", expenseHandler.ListExpenses)
			expenses.GET("/:id", expenseHandler.GetExpense)
			expenses.PUT("/:id", expenseHandler.UpdateExpense)
			expenses.DELETE("/:id", expenseHandler.DeleteExpense)
		}

		// Expense Report routes
		reports := api.Group("/reports")
		{
			reports.POST("", reportHandler.CreateReport)
			reports.GET("", reportHandler.ListReports)
			reports.GET("/:id", reportHandler.GetReport)
			reports.POST("/:id/expenses", reportHandler.AddExpensesToReport)
			reports.PUT("/:id/submit", reportHandler.SubmitReport)
		}
	}

	return router
}

