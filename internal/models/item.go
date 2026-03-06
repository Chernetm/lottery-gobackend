package models

import (
	"time"
)

type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	ImageUrl    *string   `json:"imageUrl"`
	RetailPrice float64   `json:"retailPrice"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	Lotteries   []Lottery `gorm:"foreignKey:ItemID" json:"lotteries,omitempty"`
}
