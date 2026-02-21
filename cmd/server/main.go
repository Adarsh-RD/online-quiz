package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"online-quiz/internal/domain"
	"online-quiz/internal/handler"
	"online-quiz/internal/middleware"
	"online-quiz/internal/repository"
	"online-quiz/internal/service"
	"online-quiz/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	defer logger.Sync()
	zapLog := logger.Log

	if err := godotenv.Load(); err != nil {
		zapLog.Info("No .env file found, relying on environment variables")
	}

	dbConfig := repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		zapLog.Fatal("Database connection failed", zap.Error(err))
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Services
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "super-secret"
	}
	jwtService := service.NewJWTService(jwtSecret)
	authService := service.NewAuthService(userRepo, jwtService)
	quizService := service.NewQuizService(quizRepo)
	sessionService := service.NewSessionService(sessionRepo, quizRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	quizHandler := handler.NewQuizHandler(quizService)
	sessionHandler := handler.NewSessionHandler(sessionService)

	// Router
	r := gin.Default()

	// Public Routes
	r.POST("/api/auth/register", authHandler.Register)
	r.POST("/api/auth/login", authHandler.Login)

	// Protected Routes
	api := r.Group("/api")
	api.Use(middleware.AuthorizeJWT(jwtService))
	{
		// Teacher only
		teacherOnly := api.Group("")
		teacherOnly.Use(middleware.RequireRole(domain.RoleTeacher))
		{
			teacherOnly.POST("/quizzes", quizHandler.CreateQuiz)
		}

		// Student only
		studentOnly := api.Group("")
		studentOnly.Use(middleware.RequireRole(domain.RoleStudent))
		{
			studentOnly.GET("/quizzes", quizHandler.ListQuizzes)
			studentOnly.POST("/sessions/start", sessionHandler.StartSession)
			studentOnly.POST("/sessions/answer", sessionHandler.SubmitAnswer)
			studentOnly.POST("/sessions/submit", sessionHandler.SubmitQuiz)
			studentOnly.POST("/sessions/tab-switch", sessionHandler.TabSwitchEvent)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		zapLog.Info("Server starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLog.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zapLog.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zapLog.Fatal("Server forced to shutdown", zap.Error(err))
	}

	zapLog.Info("Server exiting gracefully")
}
