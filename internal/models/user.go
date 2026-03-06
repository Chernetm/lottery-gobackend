package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
	StatusBlocked  UserStatus = "BLOCKED"
)

type User struct {
	ID            string     `gorm:"primaryKey;type:varchar(191)" json:"id"`
	Email         *string    `gorm:"uniqueIndex;type:varchar(191)" json:"email"`
	PhoneNumber   string     `gorm:"uniqueIndex;type:varchar(191)" json:"phoneNumber"`
	Password      string     `json:"password"`
	FullName      *string    `json:"fullName"`
	Role          string     `gorm:"type:varchar(20);default:'USER'" json:"role"`
	Status        UserStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	WalletBalance float64    `gorm:"default:0" json:"walletBalance"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	Coupons       []Coupon   `gorm:"foreignKey:UserID" json:"coupons,omitempty"`
	Lotteries     []Lottery  `gorm:"foreignKey:WinnerID" json:"lotteries,omitempty"`
	Tickets       []Ticket   `gorm:"foreignKey:UserID" json:"tickets,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
