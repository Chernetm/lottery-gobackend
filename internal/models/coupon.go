package models

import (
	"time"
)

type CouponType string

const (
	CouponPercentage  CouponType = "PERCENTAGE"
	CouponFixedAmount CouponType = "FIXED_AMOUNT"
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
	Type        CouponType   `gorm:"type:enum('PERCENTAGE','FIXED_AMOUNT')" json:"type"`
	Value       float64      `json:"value"`
	MinSpend    *float64     `json:"minSpend"`
	MaxDiscount *float64     `json:"maxDiscount"`
	Status      CouponStatus `gorm:"type:enum('ACTIVE','USED','EXPIRED');default:'ACTIVE'" json:"status"`
	UserID      string       `json:"userId"`
	ExpiresAt   *time.Time   `json:"expiresAt"`
	UsedAt      *time.Time   `json:"usedAt"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	User        User         `gorm:"foreignKey:UserID" json:"user"`
}
