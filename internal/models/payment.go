package models

import (
	"time"
)

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

type Payment struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	TransactionRef string        `gorm:"uniqueIndex:idx_tx_ref;type:varchar(191)" json:"transactionRef"`
	UserID         string        `json:"userId"`
	LotteryID      uint          `json:"lotteryId"`
	Quantity       int           `json:"quantity"`
	Amount         float64       `json:"amount"`
	Currency       string        `gorm:"default:'ETB'" json:"currency"`
	Status         PaymentStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	CheckoutURL    string        `json:"checkoutUrl"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updatedAt"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Lottery Lottery `gorm:"foreignKey:LotteryID" json:"-"`
}
