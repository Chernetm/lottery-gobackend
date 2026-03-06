package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminStatus string

const (
	AdminStatusActive   AdminStatus = "ACTIVE"
	AdminStatusInactive AdminStatus = "INACTIVE"
)

type Admin struct {
	ID          string      `gorm:"primaryKey;type:varchar(191)" json:"id"`
	Email       string      `gorm:"uniqueIndex;type:varchar(191)" json:"email"`
	PhoneNumber string      `gorm:"type:varchar(20)" json:"phoneNumber"`
	Password    string      `json:"-"`
	FullName    string      `json:"fullName"`
	Role        string      `gorm:"type:varchar(20);default:'ADMIN'" json:"role"`
	Status      AdminStatus `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return
}
