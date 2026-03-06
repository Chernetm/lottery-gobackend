package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
	"math/big"
	"time"
)

type AdminService struct {
	itemRepo    *repo.ItemRepo
	lotteryRepo *repo.LotteryRepo
	userRepo    *repo.UserRepo
	ticketRepo  *repo.TicketRepo
}

func NewAdminService(itemRepo *repo.ItemRepo, lotteryRepo *repo.LotteryRepo, userRepo *repo.UserRepo, ticketRepo *repo.TicketRepo) *AdminService {
	return &AdminService{
		itemRepo:    itemRepo,
		lotteryRepo: lotteryRepo,
		userRepo:    userRepo,
		ticketRepo:  ticketRepo,
	}
}

func (s *AdminService) CreateItem(name string, description *string, imageUrl *string, retailPrice float64) (*models.Item, error) {
	item := &models.Item{
		Name:        name,
		Description: description,
		ImageUrl:    imageUrl,
		RetailPrice: retailPrice,
		IsActive:    true,
	}

	if err := s.itemRepo.Create(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *AdminService) DeleteItem(id uint) error {
	return s.itemRepo.Delete(id)
}

func (s *AdminService) GetAllItems() ([]models.Item, error) {
	return s.itemRepo.FindAll()
}

func (s *AdminService) UpdateLotteryStatus(lotteryID uint, status models.LotteryStatus) error {
	lottery, err := s.lotteryRepo.FindByID(lotteryID)
	if err != nil {
		return err
	}

	lottery.Status = status
	return s.lotteryRepo.Update(lottery)
}

func (s *AdminService) UpdateLottery(id uint, ticketPrice float64, minTickets int, maxTickets *int, status models.LotteryStatus) error {
	lottery, err := s.lotteryRepo.FindByID(id)
	if err != nil {
		return err
	}

	lottery.TicketPrice = ticketPrice
	lottery.MinTickets = minTickets
	lottery.MaxTickets = maxTickets
	lottery.Status = status
	return s.lotteryRepo.Update(lottery)
}

func (s *AdminService) DrawWinner(lotteryID uint) ([]models.LotteryPrize, error) {
	lottery, err := s.lotteryRepo.FindByID(lotteryID)
	if err != nil {
		return nil, err
	}

	if lottery.Status != models.LotteryActive {
		return nil, errors.New("lottery is not active")
	}

	tickets, err := s.ticketRepo.FindActiveByLotteryID(lotteryID)
	if err != nil {
		return nil, err
	}

	if len(tickets) == 0 {
		return nil, errors.New("no active tickets for this lottery")
	}

	if len(tickets) < lottery.MinTickets {
		return nil, fmt.Errorf("minimum tickets (%d) not reached, only %d sold", lottery.MinTickets, len(tickets))
	}

	// We need to draw a winner for each prize
	if len(lottery.Prizes) == 0 {
		return nil, errors.New("no prizes defined for this lottery")
	}

	// Shuffle tickets or pick randomly of size len(prizes)
	// For simplicity, we'll pick unique winners if possible, or allow one user to win multiple prizes if they have multiple tickets.
	// The requirement is usually "unique ticket wins".

	winnerTicketIDToPrizeID := make(map[uint]uint)
	var prizes []models.LotteryPrize = lottery.Prizes

	availableTickets := make([]models.Ticket, len(tickets))
	copy(availableTickets, tickets)

	for i := range prizes {
		if len(availableTickets) == 0 {
			break
		}

		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(availableTickets))))
		winnerTicket := availableTickets[n.Int64()]

		// Correct pointer assignment (allocate new variables on heap)
		uid := winnerTicket.UserID
		prizes[i].WinnerID = &uid

		tid := winnerTicket.ID
		prizes[i].WinnerTicketID = &tid

		tnum := winnerTicket.TicketNumber
		prizes[i].WinnerTicketNumber = &tnum

		winnerTicketIDToPrizeID[winnerTicket.ID] = prizes[i].ID

		// Set legacy WinnerID for the rank 1 prize
		if prizes[i].Rank == 1 {
			lottery.WinnerID = &uid
		}

		// Remove the winning ticket from available pool
		availableTickets = append(availableTickets[:n.Int64()], availableTickets[n.Int64()+1:]...)
	}

	// Atomic update of all tickets for this lottery
	if err := s.ticketRepo.UpdateStatusesAfterMultiDraw(lotteryID, winnerTicketIDToPrizeID); err != nil {
		return nil, err
	}

	// Update prizes with winners in DB
	for _, p := range prizes {
		if err := s.lotteryRepo.UpdatePrize(&p); err != nil {
			return nil, err
		}
	}

	lottery.Status = models.LotteryDrawn
	now := time.Now()
	lottery.DrawnAt = &now

	if err := s.lotteryRepo.Update(lottery); err != nil {
		return nil, err
	}

	return prizes, nil
}

func (s *AdminService) GetStats() (*models.AdminStats, error) {
	activeLotteries, err := s.lotteryRepo.GetActiveCount()
	if err != nil {
		return nil, err
	}

	totalRevenue, err := s.ticketRepo.GetTotalRevenue()
	if err != nil {
		return nil, err
	}

	totalUsers, err := s.userRepo.GetTotalCount()
	if err != nil {
		return nil, err
	}

	return &models.AdminStats{
		TotalRevenue:    totalRevenue,
		ActiveLotteries: activeLotteries,
		TotalUsers:      totalUsers,
	}, nil
}

func (s *AdminService) GetAllTickets() ([]models.Ticket, error) {
	return s.ticketRepo.FindAll()
}

func (s *AdminService) UpdateTicketStatus(ticketID uint, status models.TicketStatus) error {
	ticket, err := s.ticketRepo.FindByID(ticketID)
	if err != nil {
		return err
	}

	ticket.Status = status
	return s.ticketRepo.Update(ticket)
}
