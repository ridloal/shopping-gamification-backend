package main

import (
	"database/sql"
	"fmt"
	"log"
	"shopping-gamification/internal/delivery/http/handler"
	"shopping-gamification/internal/repository/postgres"
	"shopping-gamification/internal/usecase"
	"shopping-gamification/pkg/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create the connection string PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the repository and usecase
	repo := postgres.NewRepository(db)
	productUsecase := usecase.NewProductUsecase(repo)
	claimUsecase := usecase.NewClaimUsecase(repo)

	// Initialize the Gin engine
	r := gin.Default()

	// Initialize handler
	handler.NewProductHandler(r, productUsecase)
	handler.NewClaimHandler(r, claimUsecase)

	r.Run(":8080")
}
