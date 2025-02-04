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
	"shopping-gamification/internal/repository/redis"
	"shopping-gamification/internal/usecase"
	"shopping-gamification/pkg/config"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	rdbDependency "github.com/redis/go-redis/v9"
)

func main() {
	cfg := loadConfig()
	db := initializeDatabase(cfg)
	rdb := initializeRedis(cfg)

	// Initialize repositories
	postgresRepo := postgres.NewRepository(db)
	var redisRepo *redis.Repository
	if rdb != nil {
		redisRepo = redis.NewRepository(rdb)
	}

	// Initialize usecases
	productUsecase := usecase.NewProductUsecase(postgresRepo, redisRepo)
	claimUsecase := usecase.NewClaimUsecase(postgresRepo, postgresRepo)
	pageUsecase := usecase.NewPageUsecase(postgresRepo, redisRepo)

	r := initializeGinEngine(db, rdb, &productUsecase, &claimUsecase, &pageUsecase)
	startServer(r, db)
}

func loadConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: %v", err)
	}
	return cfg
}

func initializeDatabase(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *sql.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
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

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db
}

func initializeRedis(cfg *config.Config) *rdbDependency.Client {
	if cfg.RedisAddress == "" {
		log.Println("Redis configuration not available, bypassing Redis")
		return nil
	}

	rdb := rdbDependency.NewClient(&rdbDependency.Options{
		Addr:     cfg.RedisAddress,
		Username: cfg.RedisUser,
		Password: cfg.RedisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Failed to connect to Redis: ", err)
		return nil
	}

	log.Println("Connected to Redis successfully")
	return rdb
}

func initializeGinEngine(db *sql.DB, rdb *rdbDependency.Client, productUsecase *usecase.ProductUsecase, claimUsecase *usecase.ClaimUsecase, pageUsecase *usecase.PageUsecase) *gin.Engine {
	r := gin.Default()

	r.GET("/health-check", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "message": "database connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/flush-redis", func(c *gin.Context) {
		err := rdb.FlushAll(context.Background()).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to flush redis"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	handler.NewProductHandler(r, *productUsecase)
	handler.NewClaimHandler(r, *claimUsecase)
	handler.NewPageHandler(r, *pageUsecase)

	return r
}

func startServer(r *gin.Engine, db *sql.DB) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	if err := db.Close(); err != nil {
		log.Fatal("Failed to close database connection: ", err)
	}

	log.Println("Server exiting")
}
