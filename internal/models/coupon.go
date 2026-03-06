package models

import (
	"time"
)

type CouponType string

const (
	CouponPercentage  CouponType   = "PERCENTAGE"
	CouponFixedAmount CouponType   = "FIXED_AMOUNT"
	CouponFreeTicket  CouponType   = "FREE_TICKET"
)

type CouponStatus string

const (
	CouponActive  CouponStatus = "ACTIVE"
	CouponUsed    CouponStatus = "USED"
	CouponExpired CouponStatus = "EXPIRED"
)

type Coupon struct {
	ID          uint         `gorm:"primaryKey" json:"id"`
	Code        string       `gorm:"uniqueIndex;size:191" json:"code"`
	Type        CouponType   `gorm:"type:varchar(20)" json:"type"`
	Value       float64      `json:"value"`
	MinSpend    *float64     `json:"minSpend"`
	MaxDiscount *float64     `json:"maxDiscount"`
	Status      CouponStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	UserID      string       `json:"userId"`
	LotteryID   *uint        `json:"lotteryId"` // Optional: specific lottery for this coupon
	ExpiresAt   *time.Time   `json:"expiresAt"`
	UsedAt      *time.Time   `json:"usedAt"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	User        User         `gorm:"foreignKey:UserID" json:"user"`
	Lottery     *Lottery     `gorm:"foreignKey:LotteryID" json:"lottery"`
}
