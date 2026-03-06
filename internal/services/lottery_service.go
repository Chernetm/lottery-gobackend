package services

import (
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
)

type LotteryService struct {
	lotteryRepo *repo.LotteryRepo
	itemRepo    *repo.ItemRepo
}

func NewLotteryService(lotteryRepo *repo.LotteryRepo, itemRepo *repo.ItemRepo) *LotteryService {
	return &LotteryService{
		lotteryRepo: lotteryRepo,
		itemRepo:    itemRepo,
	}
}

func (s *LotteryService) CreateLottery(itemID uint, ticketPrice float64, minTickets int, maxTickets *int, prizes []models.LotteryPrize) (*models.Lottery, error) {
	// If itemID is 0, we use the first prize's itemID for the legacy column
	if itemID == 0 && len(prizes) > 0 {
		itemID = prizes[0].ItemID
	}

	lottery := &models.Lottery{
		ItemID:      itemID,
		TicketPrice: ticketPrice,
		MinTickets:  minTickets,
		MaxTickets:  maxTickets,
		Status:      models.LotteryActive,
		Prizes:      prizes,
	}

	if err := s.lotteryRepo.Create(lottery); err != nil {
		return nil, err
	}

	return lottery, nil
}

func (s *LotteryService) GetAllLotteries(status string) ([]models.Lottery, error) {
	return s.lotteryRepo.FindAll(status)
}

func (s *LotteryService) GetLotteryByID(id uint) (*models.Lottery, error) {
	return s.lotteryRepo.FindByID(id)
}

func (s *LotteryService) GetLotteryItems() ([]models.Item, error) {
	return s.itemRepo.FindAll()
}
