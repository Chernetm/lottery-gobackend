package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type ItemRepo struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

func (r *ItemRepo) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *ItemRepo) FindAll() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *ItemRepo) FindByID(id uint) (*models.Item, error) {
	var item models.Item
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepo) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *ItemRepo) Delete(id uint) error {
	return r.db.Delete(&models.Item{}, id).Error
}
