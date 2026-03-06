package models

import (
	"time"
)

type LotteryStatus string

const (
	LotteryDraft     LotteryStatus = "DRAFT"
	LotteryActive    LotteryStatus = "ACTIVE"
	LotteryLocked    LotteryStatus = "LOCKED"
	LotteryDrawn     LotteryStatus = "DRAWN"
	LotteryCompleted LotteryStatus = "COMPLETED"
	LotteryCancelled LotteryStatus = "CANCELLED"
)

type Lottery struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ItemID         uint           `json:"itemId"` // Legacy: First prize item
	TicketPrice    float64        `json:"ticketPrice"`
	MinTickets     int            `json:"minTickets"`
	MaxTickets     *int           `json:"maxTickets"`
	TotalTickets   int            `gorm:"default:0" json:"totalTickets"`
	Status         LotteryStatus  `gorm:"type:enum('DRAFT','ACTIVE','LOCKED','DRAWN','CANCELLED','COMPLETED');default:'ACTIVE'" json:"status"`
	WinnerID       *string        `json:"winnerId"` // Legacy: First prize winner
	DrawAt         *time.Time     `json:"drawAt"`
	DrawnAt        *time.Time     `json:"drawnAt"`
	ServerSeedHash *string        `json:"serverSeedHash"`
	ServerSeed     *string        `json:"serverSeed"`
	PublicSeed     *string        `json:"publicSeed"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	Item           Item           `gorm:"foreignKey:ItemID" json:"item"`
	Winner         *User          `gorm:"foreignKey:WinnerID" json:"winner"`
	Tickets        []Ticket       `gorm:"foreignKey:LotteryID" json:"tickets,omitempty"`
	Prizes         []LotteryPrize `gorm:"foreignKey:LotteryID" json:"prizes,omitempty"`
}

type LotteryPrize struct {
	ID                 uint    `gorm:"primaryKey" json:"id"`
	LotteryID          uint    `json:"lotteryId"`
	ItemID             uint    `json:"itemId"`
	Rank               int     `json:"rank"` // 1 for 1st, 2 for 2nd, etc.
	WinnerID           *string `json:"winnerId"`
	WinnerTicketID     *uint   `json:"winnerTicketId"`
	WinnerTicketNumber *int    `json:"winnerTicketNumber"`
	Item               Item    `gorm:"foreignKey:ItemID" json:"item"`
	Winner             *User   `gorm:"foreignKey:WinnerID" json:"winner"`
	WinnerTicket       *Ticket `gorm:"foreignKey:WinnerTicketID" json:"winnerTicket"`
}
