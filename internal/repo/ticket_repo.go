package repo

import (
	"lottery-backend/internal/models"

	"gorm.io/gorm"
)

type TicketRepo struct {
	db *gorm.DB
}

func NewTicketRepo(db *gorm.DB) *TicketRepo {
	return &TicketRepo{db: db}
}

func (r *TicketRepo) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepo) FindByUserID(userID string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Preload("Lottery.Item").
		Preload("Lottery.Prizes").
		Preload("Lottery.Prizes.Item").
		Preload("Lottery.Prizes.Winner").
		Preload("Lottery.Prizes.WinnerTicket").
		Preload("WonPrize.Item").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepo) CountByLotteryID(lotteryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Ticket{}).Where("lottery_id = ?", lotteryID).Count(&count).Error
	return count, err
}

func (r *TicketRepo) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.db.Preload("Lottery.Item").
		Preload("Lottery.Prizes").
		Preload("Lottery.Prizes.Item").
		Preload("Lottery.Prizes.Winner").
		Preload("Lottery.Prizes.WinnerTicket").
		Preload("User").
		Preload("WonPrize.Item").
		First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepo) Update(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *TicketRepo) FindAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Preload("Lottery.Item").
		Preload("Lottery.Prizes").
		Preload("Lottery.Prizes.Item").
		Preload("Lottery.Prizes.Winner").
		Preload("Lottery.Prizes.WinnerTicket").
		Preload("User").
		Preload("WonPrize.Item").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepo) GetTotalRevenue() (float64, error) {
	var total float64
	err := r.db.Model(&models.Ticket{}).Select("COALESCE(SUM(purchase_price), 0)").Scan(&total).Error
	return total, err
}

func (r *TicketRepo) FindActiveByLotteryID(lotteryID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Where("lottery_id = ? AND status = ?", lotteryID, models.TicketActive).Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepo) UpdateStatusesAfterMultiDraw(lotteryID uint, winnerTicketIDToPrizeID map[uint]uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Set winners to WON and associate with specific prize
		for ticketID, prizeID := range winnerTicketIDToPrizeID {
			if err := tx.Model(&models.Ticket{}).Where("id = ?", ticketID).
				Updates(map[string]interface{}{
					"status":       models.TicketWon,
					"is_revealed":  false,
					"won_prize_id": prizeID,
				}).Error; err != nil {
				return err
			}
		}

		// Get all winner ticket IDs
		winnerTicketIDs := make([]uint, 0, len(winnerTicketIDToPrizeID))
		for tid := range winnerTicketIDToPrizeID {
			winnerTicketIDs = append(winnerTicketIDs, tid)
		}

		// Set all other ACTIVE tickets for this lottery to LOST
		query := tx.Model(&models.Ticket{}).Where("lottery_id = ? AND status = ?", lotteryID, models.TicketActive)
		if len(winnerTicketIDs) > 0 {
			query = query.Where("id NOT IN ?", winnerTicketIDs)
		}

		if err := query.Updates(map[string]interface{}{"status": models.TicketLost, "is_revealed": false}).Error; err != nil {
			return err
		}

		return nil
	})
}
