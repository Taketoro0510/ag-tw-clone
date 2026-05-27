package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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
	bmRepo := repository.NewBookmarkRepository(db, postRepo)
	commentRepo := repository.NewCommentRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// UseCases
	authUC := usecase.NewAuthUseCase(firebaseAuth, jwtService, userRepo, db)
	userUC := usecase.NewUserUseCase(userRepo, postRepo, likeRepo, bmRepo, followRepo)
	postUC := usecase.NewPostUseCase(postRepo, likeRepo, bmRepo, db)
	likeUC := usecase.NewLikeUseCase(likeRepo)
	bmUC := usecase.NewBookmarkUseCase(bmRepo, likeRepo)
	commentUC := usecase.NewCommentUseCase(commentRepo, userRepo, postRepo)
	followUC := usecase.NewFollowUseCase(followRepo, userRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authUC)
	userHandler := handler.NewUserHandler(userUC)
	postHandler := handler.NewPostHandler(postUC, likeUC)
	bmHandler := handler.NewBookmarkHandler(bmUC)
	commentHandler := handler.NewCommentHandler(commentUC)
	followHandler := handler.NewFollowHandler(followUC)

	e := echo.New()
	e.HTTPErrorHandler = middleware.ErrorHandler(logger)
	e.Use(middleware.Logger(logger))
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

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

		// Post rate limit
		postLimit := middleware.RateLimit(100) // 100 per min

		protected.GET("/posts", postHandler.ListPosts)
		protected.POST("/posts", postHandler.CreatePost, postLimit)
		protected.GET("/posts/:id", postHandler.GetPost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.POST("/posts/:id/likes", postHandler.LikePost, postLimit)
		protected.DELETE("/posts/:id/likes", postHandler.UnlikePost, postLimit)
		protected.POST("/posts/:id/bookmarks", bmHandler.BookmarkPost, postLimit)
		protected.DELETE("/posts/:id/bookmarks", bmHandler.UnbookmarkPost, postLimit)

		protected.GET("/posts/:id/comments", commentHandler.ListComments)
		protected.POST("/posts/:id/comments", commentHandler.CreateComment, postLimit)
		protected.DELETE("/posts/:id/comments/:commentId", commentHandler.DeleteComment)
		protected.POST("/comments/:id/likes", commentHandler.LikeComment, postLimit)
		protected.DELETE("/comments/:id/likes", commentHandler.UnlikeComment, postLimit)
		protected.POST("/comments/:id/bookmarks", commentHandler.BookmarkComment, postLimit)
		protected.DELETE("/comments/:id/bookmarks", commentHandler.UnbookmarkComment, postLimit)

		protected.GET("/users/:id", userHandler.GetUser)
		protected.GET("/users/:id/posts", userHandler.ListUserPosts)
		protected.POST("/users/:id/follow", followHandler.FollowUser, postLimit)
		protected.DELETE("/users/:id/follow", followHandler.UnfollowUser, postLimit)
		protected.GET("/users/:id/followers", followHandler.ListFollowers)
		protected.GET("/users/:id/following", followHandler.ListFollowings)

		protected.GET("/bookmarks", bmHandler.ListBookmarks)
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
