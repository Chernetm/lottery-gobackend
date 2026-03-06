package handlers

import (
	"lottery-backend/internal/models"
	"lottery-backend/internal/services"
	"net/http"
	"strconv"
	"time"

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
		LotteryID  uint   `json:"lotteryId" binding:"required"`
		Quantity   int    `json:"quantity" binding:"required"`
		CouponCode string `json:"couponCode"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userId")

	// If coupon is provided, try to process it first
	if input.CouponCode != "" {
		// For free tickets via coupon, we currently support quantity 1
		if input.Quantity != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Coupons currently support only one free ticket at a time"})
			return
		}

		// Generate random ticket number
		ticketNumStr := strconv.FormatInt(time.Now().UnixNano(), 10)
		if len(ticketNumStr) > 9 {
			ticketNumStr = ticketNumStr[len(ticketNumStr)-9:]
		}
		ticketNumber, _ := strconv.Atoi(ticketNumStr)

		ticket, err := h.ticketService.PurchaseTicket(userID, input.LotteryID, ticketNumber, input.CouponCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Ticket purchased successfully with coupon",
			"ticket":  ticket,
		})
		return
	}

	userEmail := c.GetString("userEmail")
	fullName := c.GetString("fullName")

	var emailPtr *string
	if userEmail != "" {
		emailPtr = &userEmail
	}

	user := &models.User{
		ID:       userID,
		Email:    emailPtr,
		FullName: &fullName,
	}

	lottery, err := h.ticketService.GetLotteryByID(input.LotteryID)
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
