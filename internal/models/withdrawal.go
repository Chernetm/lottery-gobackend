package models

import (
	"time"
)

type WithdrawalStatus string

const (
	WithdrawalPending  WithdrawalStatus = "PENDING"
	WithdrawalApproved WithdrawalStatus = "APPROVED"
	WithdrawalRejected WithdrawalStatus = "REJECTED"
)

type Withdrawal struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	UserID    string           `gorm:"type:varchar(191);index" json:"userId"`
	Amount    float64          `json:"amount"`
	Status    WithdrawalStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updatedAt"`
	User      User             `gorm:"foreignKey:UserID" json:"user"`
}
