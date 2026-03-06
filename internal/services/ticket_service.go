package services

import (
	"errors"
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
	"time"
)

type TicketService struct {
	ticketRepo  *repo.TicketRepo
	userRepo    *repo.UserRepo
	lotteryRepo *repo.LotteryRepo
	couponRepo  *repo.CouponRepo
}

func NewTicketService(ticketRepo *repo.TicketRepo, userRepo *repo.UserRepo, lotteryRepo *repo.LotteryRepo, couponRepo *repo.CouponRepo) *TicketService {
	return &TicketService{
		ticketRepo:  ticketRepo,
		userRepo:    userRepo,
		lotteryRepo: lotteryRepo,
		couponRepo:  couponRepo,
	}
}

func (s *TicketService) PurchaseTicket(userID string, lotteryID uint, ticketNumber int, couponCode string) (*models.Ticket, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	lottery, err := s.lotteryRepo.FindByID(lotteryID)
	if err != nil {
		return nil, err
	}

	isFree := false
	var coupon *models.Coupon

	if couponCode != "" {
		coupon, err = s.couponRepo.FindActiveByCodeAndUser(couponCode, userID)
		if err != nil {
			return nil, errors.New("invalid or expired coupon")
		}

		if coupon.Type != models.CouponFreeTicket {
			return nil, errors.New("this coupon is not for a free ticket")
		}

		if coupon.LotteryID == nil || *coupon.LotteryID != lotteryID {
			return nil, errors.New("this coupon is not valid for this lottery")
		}

		if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("coupon has expired")
		}

		isFree = true
	}

	if !isFree {
		if user.WalletBalance < lottery.TicketPrice {
			return nil, errors.New("insufficient balance")
		}
		user.WalletBalance -= lottery.TicketPrice
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	} else {
		// Mark coupon as used
		now := time.Now()
		coupon.Status = models.CouponUsed
		coupon.UsedAt = &now
		if err := s.couponRepo.Update(coupon); err != nil {
			return nil, err
		}
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
