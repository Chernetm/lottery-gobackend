package handlers

import (
	"lottery-backend/internal/models"
	"lottery-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WithdrawalHandler struct {
	withdrawalService *services.WithdrawalService
}

func NewWithdrawalHandler(withdrawalService *services.WithdrawalService) *WithdrawalHandler {
	return &WithdrawalHandler{withdrawalService: withdrawalService}
}

func (h *WithdrawalHandler) RequestWithdrawal(c *gin.Context) {
	userID := c.GetString("userId")
	var input struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	withdrawal, err := h.withdrawalService.RequestWithdrawal(userID, input.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawal)
}

func (h *WithdrawalHandler) GetMyWithdrawals(c *gin.Context) {
	userID := c.GetString("userId")
	withdrawals, err := h.withdrawalService.GetUserWithdrawals(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

func (h *WithdrawalHandler) GetAllWithdrawals(c *gin.Context) {
	withdrawals, err := h.withdrawalService.GetAllWithdrawals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

func (h *WithdrawalHandler) UpdateWithdrawalStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid withdrawal id"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.withdrawalService.UpdateWithdrawalStatus(uint(id), models.WithdrawalStatus(input.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "withdrawal status updated"})
}
