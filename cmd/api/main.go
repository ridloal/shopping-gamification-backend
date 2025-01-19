package main

import (
	"database/sql"
	"log"
	"shopping-gamification/internal/delivery/http/handler"
	"shopping-gamification/internal/repository/mysql"
	"shopping-gamification/internal/usecase"
	"shopping-gamification/pkg/config"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create the connection string
	dsn := cfg.DBUser + ":" + cfg.DBPassword + "@tcp(" + cfg.DBHost + ":" + cfg.DBPort + ")/" + cfg.DBName + "?parseTime=true"

	// Connect to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the repository and usecase
	repo := mysql.NewRepository(db)
	productUsecase := usecase.NewProductUsecase(repo)
	claimUsecase := usecase.NewClaimUsecase(repo)

	// Initialize the Gin engine
	r := gin.Default()

	// Initialize handler
	handler.NewProductHandler(r, productUsecase)
	handler.NewClaimHandler(r, claimUsecase)

	r.Run(":8080")
}
