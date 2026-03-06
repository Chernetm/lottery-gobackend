package models

import (
	"time"
)

type TicketStatus string

const (
	TicketActive   TicketStatus = "ACTIVE"
	TicketWon      TicketStatus = "WON"
	TicketLost     TicketStatus = "LOST"
	TicketRefunded TicketStatus = "REFUNDED"
)

type Ticket struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	TicketNumber  int          `gorm:"uniqueIndex:idx_lottery_ticket" json:"ticketNumber"`
	UserID        string       `json:"userId"`
	LotteryID     uint         `gorm:"uniqueIndex:idx_lottery_ticket" json:"lotteryId"`
	Status        TicketStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	IsRevealed    bool         `gorm:"default:false" json:"isRevealed"`
	PurchasePrice float64      `json:"purchasePrice"`
	MessageRead   bool
	CreatedAt     time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	WonPrizeID    *uint         `json:"wonPrizeId"`
	Lottery       Lottery       `gorm:"foreignKey:LotteryID" json:"lottery"`
	User          User          `gorm:"foreignKey:UserID" json:"user"`
	WonPrize      *LotteryPrize `gorm:"foreignKey:WonPrizeID" json:"wonPrize"`
}
