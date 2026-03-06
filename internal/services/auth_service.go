package services

import (
	"errors"
	"lottery-backend/internal/config"
	"lottery-backend/internal/models"
	"lottery-backend/internal/repo"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repo.UserRepo
	adminRepo *repo.AdminRepo
}

func NewAuthService(userRepo *repo.UserRepo, adminRepo *repo.AdminRepo) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		adminRepo: adminRepo,
	}
}

func (s *AuthService) Register(email *string, password, phoneNumber, fullName string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:          uuid.New().String(),
		Email:       email,
		Password:    string(hashedPassword),
		PhoneNumber: phoneNumber,
		FullName:    &fullName,
		Role:        "USER",
		Status:      models.StatusActive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(identifier, password string) (string, *models.User, error) {
	var user *models.User
	var err error

	// Check if identifier is phone number (starts with 0) or email
	if len(identifier) > 0 && identifier[0] == '0' {
		// Format phone number: replace 0 with 251
		formattedPhone := "251" + identifier[1:]
		user, err = s.userRepo.FindByPhoneNumber(formattedPhone)
	} else {
		user, err = s.userRepo.FindByEmail(identifier)
	}

	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"type": "user",
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *AuthService) AdminRegister(email, password, fullName, phoneNumber string) (*models.Admin, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &models.Admin{
		Email:       email,
		Password:    string(hashedPassword),
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Role:        "ADMIN",
		Status:      models.AdminStatusActive,
	}

	if err := s.adminRepo.Create(admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func (s *AuthService) AdminLogin(identifier, password string) (string, *models.Admin, error) {
	var admin *models.Admin
	var err error

	// Check if identifier is phone number (starts with 0) or email
	if len(identifier) > 0 && identifier[0] == '0' {
		// Format phone number: replace 0 with 251
		formattedPhone := "251" + identifier[1:]
		admin, err = s.adminRepo.FindByPhoneNumber(formattedPhone)
	} else {
		admin, err = s.adminRepo.FindByEmail(identifier)
	}

	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  admin.ID,
		"type": "admin",
		"role": admin.Role,
		"exp":  time.Now().Add(time.Hour * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, admin, nil
}

func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) GetAdminByID(id string) (*models.Admin, error) {
	return s.adminRepo.FindByID(id)
}

func (s *AuthService) ChangePassword(id string, isAdmin bool, oldPassword, newPassword string) error {
	var currentHashedPassword string
	var user *models.User
	var admin *models.Admin
	var err error

	if isAdmin {
		admin, err = s.adminRepo.FindByID(id)
		if err != nil {
			return err
		}
		currentHashedPassword = admin.Password
	} else {
		user, err = s.userRepo.FindByID(id)
		if err != nil {
			return err
		}
		currentHashedPassword = user.Password
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHashedPassword), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if isAdmin {
		admin.Password = string(newHashedPassword)
		return s.adminRepo.Update(admin)
	} else {
		user.Password = string(newHashedPassword)
		return s.userRepo.Update(user)
	}
}
