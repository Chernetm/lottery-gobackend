package services

import (
	"errors"
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
)

type WithdrawalService struct {
	withdrawalRepo *repo.WithdrawalRepo
	userRepo       *repo.UserRepo
}

func NewWithdrawalService(withdrawalRepo *repo.WithdrawalRepo, userRepo *repo.UserRepo) *WithdrawalService {
	return &WithdrawalService{
		withdrawalRepo: withdrawalRepo,
		userRepo:       userRepo,
	}
}

func (s *WithdrawalService) RequestWithdrawal(userID string, amount float64) (*models.Withdrawal, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.WalletBalance < amount {
		return nil, errors.New("insufficient wallet balance")
	}

	// Deduct from wallet immediately? Or wait for approval?
	// Usually, it's safer to deduct immediately to "reserve" the funds.
	user.WalletBalance -= amount
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	withdrawal := &models.Withdrawal{
		UserID: userID,
		Amount: amount,
		Status: models.WithdrawalPending,
	}

	if err := s.withdrawalRepo.Create(withdrawal); err != nil {
		return nil, err
	}

	return withdrawal, nil
}

func (s *WithdrawalService) GetUserWithdrawals(userID string) ([]models.Withdrawal, error) {
	return s.withdrawalRepo.FindByUserID(userID)
}

func (s *WithdrawalService) GetAllWithdrawals() ([]models.Withdrawal, error) {
	return s.withdrawalRepo.FindAll()
}

func (s *WithdrawalService) UpdateWithdrawalStatus(id uint, status models.WithdrawalStatus) error {
	withdrawal, err := s.withdrawalRepo.FindByID(id)
	if err != nil {
		return err
	}

	if withdrawal.Status != models.WithdrawalPending {
		return errors.New("can only update pending withdrawals")
	}

	if status == models.WithdrawalRejected {
		// Refund user if rejected
		user, err := s.userRepo.FindByID(withdrawal.UserID)
		if err == nil {
			user.WalletBalance += withdrawal.Amount
			s.userRepo.Update(user)
		}
	}

	withdrawal.Status = status
	return s.withdrawalRepo.Update(withdrawal)
}
