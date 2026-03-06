package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepo) FindByTransactionRef(ref string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, "transaction_ref = ?", ref).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepo) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}
