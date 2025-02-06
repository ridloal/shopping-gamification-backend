package handler

import (
	"net/http"
	"shopping-gamification/internal/delivery/http/middleware"
	"shopping-gamification/internal/domain"
	"shopping-gamification/internal/usecase"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ClaimHandler struct {
	usecase usecase.ClaimUsecase
}

func NewClaimHandler(r *gin.Engine, u usecase.ClaimUsecase, rdb *redis.Client) {
	handler := &ClaimHandler{usecase: u}

	r.POST("/claims", middleware.RateLimiter(rdb, 3, 2*time.Minute), middleware.ValidateRequest(&domain.ClaimRequestInput{}), handler.CreateClaimRequest)
	r.GET("/claims/:id", handler.GetClaimRequestByID)
	r.PATCH("/claims/:id/prizes/:prize_id", handler.UpdateClaimRequestPrize)
	r.GET("/claims/code/:code", handler.GetClaimRequestByClaimCode)
	r.GET("/claims/get-prize/:code", handler.ClaimPrize)

}

func (h *ClaimHandler) CreateClaimRequest(c *gin.Context) {
	validatedInput, exists := c.Get("validated_input")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get validated input"})
		return
	}

	req := validatedInput.(*domain.ClaimRequestInput)
	claimReq, err := h.usecase.CreateClaimRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, claimReq)
}

func (h *ClaimHandler) GetClaimRequestByID(c *gin.Context) {
	claimID := c.Param("id")
	id, err := strconv.ParseInt(claimID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claim ID"})
		return
	}
	claimReq, err := h.usecase.GetClaimRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, claimReq)
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

func (h *ClaimHandler) GetClaimRequestByClaimCode(c *gin.Context) {
	claimCode := c.Param("code")
	claimReq, err := h.usecase.GetClaimRequestByClaimCode(claimCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, claimReq)
}

func (h *ClaimHandler) ClaimPrize(c *gin.Context) {
	claimCode := c.Param("code")
	prizeReward, err := h.usecase.ClaimPrize(claimCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prizeReward)
}
