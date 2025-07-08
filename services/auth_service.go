package services

import (
	"case_study_api/dto"
	"case_study_api/entities"
	"case_study_api/repositories"
	"case_study_api/utils"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("email is already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user entity
	user := entities.User{
		Name:     req.Name,
		Email:    strings.ToLower(req.Email),
		Password: string(hashed),
		Role:     "user",
	}

	// Save user
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}

func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil || user.ID == 0 {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.AuthResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}, nil
}
