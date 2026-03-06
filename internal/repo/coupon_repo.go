package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type CouponRepo struct {
	db *gorm.DB
}

func NewCouponRepo(db *gorm.DB) *CouponRepo {
	return &CouponRepo{db: db}
}

func (r *CouponRepo) FindByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	if err := r.db.Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *CouponRepo) Update(coupon *models.Coupon) error {
	return r.db.Save(coupon).Error
}
