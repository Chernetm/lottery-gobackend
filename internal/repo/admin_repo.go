package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type AdminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AdminRepo) FindByID(id string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.First(&admin, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) FindByEmail(email string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.First(&admin, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) FindByPhoneNumber(phone string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.First(&admin, "phone_number = ?", phone).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepo) Update(admin *models.Admin) error {
	return r.db.Save(admin).Error
}
