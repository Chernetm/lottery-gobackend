package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type WithdrawalRepo struct {
	db *gorm.DB
}

func NewWithdrawalRepo(db *gorm.DB) *WithdrawalRepo {
	return &WithdrawalRepo{db: db}
}

func (r *WithdrawalRepo) Create(withdrawal *models.Withdrawal) error {
	return r.db.Create(withdrawal).Error
}

func (r *WithdrawalRepo) FindByUserID(userID string) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&withdrawals).Error
	return withdrawals, err
}

func (r *WithdrawalRepo) FindAll() ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal
	err := r.db.Preload("User").Order("created_at desc").Find(&withdrawals).Error
	return withdrawals, err
}

func (r *WithdrawalRepo) FindByID(id uint) (*models.Withdrawal, error) {
	var withdrawal models.Withdrawal
	if err := r.db.First(&withdrawal, id).Error; err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

func (r *WithdrawalRepo) Update(withdrawal *models.Withdrawal) error {
	return r.db.Save(withdrawal).Error
}
