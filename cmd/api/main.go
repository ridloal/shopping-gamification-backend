package main

import (
	"database/sql"
	"log"
	"shopping-gamification/internal/domain"
	repository "shopping-gamification/internal/repository/mysql"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/shopping_gamification?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	r := gin.Default()

	r.GET("/products", func(c *gin.Context) {
		products, err := repo.GetProducts()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, products)
	})

	r.GET("/products/:id/prize-groups", func(c *gin.Context) {
		productID := c.Param("id")
		id, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid product ID"})
			return
		}
		prizeGroups, err := repo.GetPrizeGroupsByProductID(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, prizeGroups)
	})

	r.POST("/claim-requests", func(c *gin.Context) {
		var req domain.ClaimRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := repo.CreateClaimRequest(&req); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, req)
	})

	r.PATCH("/claim-requests/:id/prize", func(c *gin.Context) {
		claimID := c.Param("id")
		id, err := strconv.ParseInt(claimID, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid claim ID"})
			return
		}
		var req struct {
			PrizeID int64 `json:"prize_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := repo.UpdateClaimRequestPrize(id, req.PrizeID); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Prize updated successfully"})
	})

	r.Run(":8080")
}
