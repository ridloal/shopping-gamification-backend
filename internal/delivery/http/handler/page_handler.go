package handler

import (
	"net/http"
	"shopping-gamification/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	usecase usecase.PageUsecase
}

func NewPageHandler(r *gin.Engine, u usecase.PageUsecase) {
	handler := &PageHandler{usecase: u}
	r.GET("/page/home", handler.GetPageHome)
}

func (h *PageHandler) GetPageHome(c *gin.Context) {
	page, err := h.usecase.GetPageHome()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, page)
}
