package handler

import (
	"net/http"
	"shopping-gamification/internal/domain"
	"shopping-gamification/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClaimHandler struct {
	usecase usecase.ClaimUsecase
}

func NewClaimHandler(r *gin.Engine, u usecase.ClaimUsecase) {
	handler := &ClaimHandler{usecase: u}
	r.POST("/claims", handler.CreateClaimRequest)
	r.PATCH("/claims/:id/prizes/:prize_id", handler.UpdateClaimRequestPrize)
}

func (h *ClaimHandler) CreateClaimRequest(c *gin.Context) {
	var req domain.ClaimRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.usecase.CreateClaimRequest(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *ClaimHandler) UpdateClaimRequestPrize(c *gin.Context) {
	claimID := c.Param("id")
	prizeID := c.Param("prize_id")
	id, err := strconv.ParseInt(claimID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}
	idPrize, err := strconv.ParseInt(prizeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prize ID"})
		return
	}
	err = h.usecase.UpdateClaimRequestPrize(id, idPrize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Claim request updated"})
}
