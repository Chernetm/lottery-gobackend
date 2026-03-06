package handlers

import (
	"lottery-backend/internal/models"
	"lottery-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketService  *services.TicketService
	paymentService *services.PaymentService
}

func NewTicketHandler(ticketService *services.TicketService, paymentService *services.PaymentService) *TicketHandler {
	return &TicketHandler{
		ticketService:  ticketService,
		paymentService: paymentService,
	}
}

func (h *TicketHandler) PurchaseTicket(c *gin.Context) {
	var input struct {
		LotteryID uint `json:"lotteryId" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userId")
	userEmail := c.GetString("userEmail") // Ensure auth middleware provides this
	fullName := c.GetString("fullName")   // Ensure auth middleware provides this

	user := &models.User{
		ID:       userID,
		Email:    userEmail,
		FullName: &fullName,
	}

	lotteryID := input.LotteryID

	// Fetch lottery to get price
	lottery, err := h.ticketService.GetLotteryByID(lotteryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lottery not found"})
		return
	}

	payment, err := h.paymentService.InitializePayment(user, lottery, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkoutUrl": payment.CheckoutURL,
		"tx_ref":      payment.TransactionRef,
	})
}

func (h *TicketHandler) GetUserTickets(c *gin.Context) {
	userID := c.GetString("userId")
	tickets, err := h.ticketService.GetUserTickets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}
func (h *TicketHandler) RevealTicket(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	userID := c.GetString("userId")
	err = h.ticketService.RevealTicket(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket revealed"})
}
