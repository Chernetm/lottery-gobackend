package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type LotteryRepo struct {
	db *gorm.DB
}

func NewLotteryRepo(db *gorm.DB) *LotteryRepo {
	return &LotteryRepo{db: db}
}

func (r *LotteryRepo) Create(lottery *models.Lottery) error {
	return r.db.Create(lottery).Error
}

func (r *LotteryRepo) FindAll(status string) ([]models.Lottery, error) {
	var lotteries []models.Lottery
	query := r.db.Preload("Item").
		Preload("Prizes").
		Preload("Prizes.Item").
		Preload("Prizes.Winner").
		Preload("Prizes.WinnerTicket").
		Preload("Winner")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&lotteries).Error
	return lotteries, err
}

func (r *LotteryRepo) FindByID(id uint) (*models.Lottery, error) {
	var lottery models.Lottery
	if err := r.db.Preload("Item").
		Preload("Prizes").
		Preload("Prizes.Item").
		Preload("Prizes.Winner").
		Preload("Prizes.WinnerTicket").
		Preload("Winner").
		First(&lottery, id).Error; err != nil {
		return nil, err
	}
	return &lottery, nil
}

func (r *LotteryRepo) Update(lottery *models.Lottery) error {
	return r.db.Save(lottery).Error
}

func (r *LotteryRepo) UpdatePrize(prize *models.LotteryPrize) error {
	return r.db.Save(prize).Error
}

func (r *LotteryRepo) GetActiveCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Lottery{}).Where("status = ?", models.LotteryActive).Count(&count).Error
	return count, err
}

func (r *LotteryRepo) GetTotalRevenue() (float64, error) {
	var total float64
	err := r.db.Model(&models.Lottery{}).Select("COALESCE(SUM(total_tickets * ticket_price), 0)").Scan(&total).Error
	return total, err
}
