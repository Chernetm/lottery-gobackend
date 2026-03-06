package handlers

import (
	"lottery-backend/internal/models"
	"lottery-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LotteryHandler struct {
	lotteryService *services.LotteryService
}

func NewLotteryHandler(lotteryService *services.LotteryService) *LotteryHandler {
	return &LotteryHandler{lotteryService: lotteryService}
}

func (h *LotteryHandler) GetAllLotteries(c *gin.Context) {
	status := c.Query("status")
	lotteries, err := h.lotteryService.GetAllLotteries(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lotteries)
}

func (h *LotteryHandler) GetLotteryByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lottery id"})
		return
	}

	lottery, err := h.lotteryService.GetLotteryByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lottery not found"})
		return
	}
	c.JSON(http.StatusOK, lottery)
}

func (h *LotteryHandler) CreateLottery(c *gin.Context) {
	var input struct {
		ItemID      uint    `json:"itemId"` // Legacy
		TicketPrice float64 `json:"ticketPrice" binding:"required"`
		MinTickets  int     `json:"minTickets" binding:"required"`
		MaxTickets  *int    `json:"maxTickets"`
		Prizes      []struct {
			ItemID uint `json:"itemId" binding:"required"`
			Rank   int  `json:"rank" binding:"required"`
		} `json:"prizes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Transform prizes or handle legacy
	var prizes []models.LotteryPrize
	if len(input.Prizes) > 0 {
		for _, p := range input.Prizes {
			prizes = append(prizes, models.LotteryPrize{
				ItemID: p.ItemID,
				Rank:   p.Rank,
			})
		}
	} else if input.ItemID != 0 {
		prizes = append(prizes, models.LotteryPrize{
			ItemID: input.ItemID,
			Rank:   1,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one prize is required"})
		return
	}

	lottery, err := h.lotteryService.CreateLottery(input.ItemID, input.TicketPrice, input.MinTickets, input.MaxTickets, prizes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lottery)
}

func (h *LotteryHandler) GetLotteryItems(c *gin.Context) {
	items, err := h.lotteryService.GetLotteryItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
