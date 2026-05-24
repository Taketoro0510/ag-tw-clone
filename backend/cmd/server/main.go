package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/koitake1/cloudcode-sns/backend/docs"
	"github.com/koitake1/cloudcode-sns/backend/internal/auth"
	"github.com/koitake1/cloudcode-sns/backend/internal/handler"
	"github.com/koitake1/cloudcode-sns/backend/internal/middleware"
	"github.com/koitake1/cloudcode-sns/backend/internal/repository"
	"github.com/koitake1/cloudcode-sns/backend/internal/usecase"
)

// @title           CloudCode SNS API
// @version         1.0
// @description     This is a sample server for SNS App MVP.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// DB Setup
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://app:app@db:5432/app?sslmode=disable"
	}

	// Migrations
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Error("failed to init migrate", "err", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to run migrate up", "err", err)
		os.Exit(1)
	}
	logger.Info("migrations applied")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect database", "err", err)
		os.Exit(1)
	}

	// Firebase Setup
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = "cloudcode-sns-local"
	}
	firebaseAuth, err := auth.NewFirebaseAuth(context.Background(), projectID)
	if err != nil {
		logger.Error("failed to init firebase auth", "err", err)
		os.Exit(1)
	}

	// JWT Setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-do-not-use-in-prod-xxxxxxxxxxxxxxxx"
	}
	jwtService := auth.NewJWTService(jwtSecret)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	likeRepo := repository.NewLikeRepository(db)

	// UseCases
	authUC := usecase.NewAuthUseCase(firebaseAuth, jwtService, userRepo, db)
	userUC := usecase.NewUserUseCase(userRepo, postRepo, likeRepo)
	postUC := usecase.NewPostUseCase(postRepo, likeRepo, db)
	likeUC := usecase.NewLikeUseCase(likeRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC)
	userHandler := handler.NewUserHandler(userUC)
	postHandler := handler.NewPostHandler(postUC, likeUC)

	e := echo.New()
	e.HTTPErrorHandler = middleware.ErrorHandler(logger)
	e.Use(middleware.Logger(logger))

	// Swagger setup
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	{
		v1.GET("/healthz", handler.Healthz)
		v1.POST("/auth/sessions", authHandler.CreateSession)
		
		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.Auth(jwtService))
		
		protected.GET("/me", userHandler.GetMe)
		protected.GET("/users/:id", userHandler.GetUser)
		protected.GET("/users/:id/posts", userHandler.ListUserPosts)
		
		// Post rate limit
		postLimit := middleware.RateLimit(5) // 5 per min
		likeLimit := middleware.RateLimit(30) // 30 per min
		
		protected.GET("/posts", postHandler.ListPosts)
		protected.POST("/posts", postHandler.CreatePost, postLimit)
		protected.GET("/posts/:id", postHandler.GetPost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.POST("/posts/:id/likes", postHandler.LikePost, likeLimit)
		protected.DELETE("/posts/:id/likes", postHandler.UnlikePost, likeLimit)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server", "port", port)
	if err := e.Start(":" + port); err != nil {
		logger.Error("server failed", "err", err)
	}
}
