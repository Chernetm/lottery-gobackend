package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lottery-backend/internal/config"
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
	"net/http"
	"time"
)

type PaymentService struct {
	paymentRepo *repo.PaymentRepo
	userRepo    *repo.UserRepo
	ticketRepo  *repo.TicketRepo
	lotteryRepo *repo.LotteryRepo
}

func NewPaymentService(
	paymentRepo *repo.PaymentRepo,
	userRepo *repo.UserRepo,
	ticketRepo *repo.TicketRepo,
	lotteryRepo *repo.LotteryRepo,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
		ticketRepo:  ticketRepo,
		lotteryRepo: lotteryRepo,
	}
}

type ChapaInitRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Email       string  `json:"email"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	TxRef       string  `json:"tx_ref"`
	CallbackURL string  `json:"callback_url"`
	ReturnURL   string  `json:"return_url"`
}

type ChapaInitResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

type ChapaVerifyResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    struct {
		Status string  `json:"status"`
		TxRef  string  `json:"tx_ref"`
		Amount float64 `json:"amount"`
	} `json:"data"`
}

func (s *PaymentService) InitializePayment(user *models.User, lottery *models.Lottery, quantity int) (*models.Payment, error) {
	// tx_ref must be < 50 chars and valid characters. UUID is ~36 chars.
	// Using a shorter hash or prefix to ensure we stay under 50.
	txRef := fmt.Sprintf("TX-%d-%s-%d", lottery.ID, user.ID[:8], time.Now().Unix())

	firstName := "User"
	if user.FullName != nil && *user.FullName != "" {
		firstName = *user.FullName
	}

	reqBody := ChapaInitRequest{
		Amount:      lottery.TicketPrice * float64(quantity),
		Currency:    "ETB",
		Email:       user.Email,
		FirstName:   firstName,
		LastName:    "Customer",
		TxRef:       txRef,
		CallbackURL: fmt.Sprintf("%s/api/payments/webhook", config.AppConfig.BaseURL),
		ReturnURL:   fmt.Sprintf("%s/payment-success?tx_ref=%s", config.AppConfig.FrontendURL, txRef),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.AppConfig.ChapaSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chapa initialization failed: %s", string(body))
	}

	var initResp ChapaInitResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return nil, err
	}

	payment := &models.Payment{
		TransactionRef: txRef,
		UserID:         user.ID,
		LotteryID:      lottery.ID,
		Quantity:       quantity,
		Amount:         lottery.TicketPrice * float64(quantity),
		Status:         models.PaymentPending,
		CheckoutURL:    initResp.Data.CheckoutURL,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) VerifyPayment(txRef string) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.chapa.co/v1/transaction/verify/%s", txRef), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.AppConfig.ChapaSecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to verify transaction with chapa")
	}

	var verifyResp ChapaVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return err
	}

	if verifyResp.Data.Status != "success" {
		return fmt.Errorf("transaction not successful: %s", verifyResp.Data.Status)
	}

	return s.FinalizePayment(txRef)
}

func (s *PaymentService) FinalizePayment(txRef string) error {
	payment, err := s.paymentRepo.FindByTransactionRef(txRef)
	if err != nil {
		return err
	}

	if payment.Status == models.PaymentSuccess {
		return nil // Already processed
	}

	payment.Status = models.PaymentSuccess
	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	// Create the tickets
	for i := 0; i < payment.Quantity; i++ {
		// In a real system, you'd want better ticket number generation
		// For now, random number that doesn't conflict easily
		ticketNumber := int(time.Now().UnixNano() % 1000000)

		ticket := &models.Ticket{
			UserID:        payment.UserID,
			LotteryID:     payment.LotteryID,
			TicketNumber:  ticketNumber,
			PurchasePrice: payment.Amount / float64(payment.Quantity),
			Status:        models.TicketActive,
		}

		if err := s.ticketRepo.Create(ticket); err != nil {
			// Log error but continue or handle as needed
			fmt.Printf("Error creating ticket %d: %v\n", i, err)
		}
	}

	// Update lottery ticket count
	lottery, err := s.lotteryRepo.FindByID(payment.LotteryID)
	if err == nil {
		lottery.TotalTickets += payment.Quantity
		s.lotteryRepo.Update(lottery)
	}

	return nil
}
