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

func (r *CouponRepo) Create(coupon *models.Coupon) error {
	return r.db.Create(coupon).Error
}

func (r *CouponRepo) FindByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	if err := r.db.Where("code = ?", code).First(&coupon).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *CouponRepo) FindActiveByCodeAndUser(code, userID string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.db.Where("code = ? AND user_id = ? AND status = ?", code, userID, models.CouponActive).
		Preload("Lottery").
		First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *CouponRepo) Update(coupon *models.Coupon) error {
	return r.db.Save(coupon).Error
}
