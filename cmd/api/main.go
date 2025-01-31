package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shopping-gamification/internal/delivery/http/handler"
	"shopping-gamification/internal/repository/postgres"
	"shopping-gamification/internal/usecase"
	"shopping-gamification/pkg/config"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: %v", err)
	}

	// Create the connection string PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Connect to the database with retry mechanism
	var db *sql.DB
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			// Test the connection
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Failed to connect to database, attempt %d/%d: %v", i+1, maxRetries, err)
		time.Sleep(time.Second * 5)
	}
	if err != nil {
		log.Fatal("Failed to connect to database after retries: ", err)
	}
	defer db.Close()

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// Initialize the repository and usecase
	repo := postgres.NewRepository(db)
	productUsecase := usecase.NewProductUsecase(repo)
	claimUsecase := usecase.NewClaimUsecase(repo)
	pageUsecase := usecase.NewPageUsecase(repo)

	// Initialize the Gin engine
	r := gin.Default()

	// Add health check endpoint
	r.GET("/health-check", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "database connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize handler
	handler.NewProductHandler(r, productUsecase)
	handler.NewClaimHandler(r, claimUsecase)
	handler.NewPageHandler(r, pageUsecase)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create server with timeout configurations
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
