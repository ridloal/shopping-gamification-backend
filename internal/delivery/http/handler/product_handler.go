package handler

import (
	"net/http"
	"shopping-gamification/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(r *gin.Engine, u usecase.ProductUsecase) {
	handler := &ProductHandler{usecase: u}
	r.GET("/products", handler.GetProducts)
	r.GET("/products/:id", handler.GetProductByID)
	r.GET("/products/:id/prize-groups", handler.GetPrizeGroupsByProductID)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.usecase.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	product, err := h.usecase.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetPrizeGroupsByProductID(c *gin.Context) {
	productID := c.Param("id")
	id, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	prizeGroups, err := h.usecase.GetPrizeGroupsByProductID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prizeGroups)
}
