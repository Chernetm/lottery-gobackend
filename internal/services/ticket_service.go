package services

import (
	"errors"
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
)

type TicketService struct {
	ticketRepo  *repo.TicketRepo
	userRepo    *repo.UserRepo
	lotteryRepo *repo.LotteryRepo
}

func NewTicketService(ticketRepo *repo.TicketRepo, userRepo *repo.UserRepo, lotteryRepo *repo.LotteryRepo) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		userRepo:    userRepo,
		lotteryRepo: lotteryRepo,
	}
}

func (s *TicketService) PurchaseTicket(userID string, lotteryID uint, ticketNumber int) (*models.Ticket, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	lottery, err := s.lotteryRepo.FindByID(lotteryID)
	if err != nil {
		return nil, err
	}

	if user.WalletBalance < lottery.TicketPrice {
		return nil, errors.New("insufficient balance")
	}

	// Transactional logic would be better here via repo, but keeping it simple for now
	user.WalletBalance -= lottery.TicketPrice
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	ticket := &models.Ticket{
		UserID:        user.ID,
		LotteryID:     lottery.ID,
		TicketNumber:  ticketNumber,
		PurchasePrice: lottery.TicketPrice,
		Status:        models.TicketActive,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *TicketService) GetUserTickets(userID string) ([]models.Ticket, error) {
	return s.ticketRepo.FindByUserID(userID)
}

func (s *TicketService) GetLotteryByID(id uint) (*models.Lottery, error) {
	return s.lotteryRepo.FindByID(id)
}

func (s *TicketService) RevealTicket(ticketID uint, userID string) error {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return err
	}

	if ticket.UserID != userID {
		return errors.New("unauthorized")
	}

	if ticket.Lottery.Status != models.LotteryDrawn {
		return errors.New("lottery not drawn yet")
	}

	ticket.IsRevealed = true
	return s.ticketRepo.Update(ticket)
}
