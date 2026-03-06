package handlers

import (
	"lottery-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// In production, verify the webhook signature using config.AppConfig.ChapaWebhookSecret
	// For now, we'll just process it

	var payload struct {
		TxRef  string `json:"tx_ref"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Status == "success" {
		if err := h.paymentService.FinalizePayment(payload.TxRef); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PaymentHandler) VerifyTransaction(c *gin.Context) {
	txRef := c.Param("tx_ref")
	if txRef == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tx_ref is required"})
		return
	}

	if err := h.paymentService.VerifyPayment(txRef); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment verified and ticket created"})
}
