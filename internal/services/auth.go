package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/group14000/golang-todo/internal/database"
	"github.com/group14000/golang-todo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      database.UserRepository
	jwtSecret string
}

func NewAuthService(repo database.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) SignUp(ctx context.Context, name, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Save to database
	return s.repo.CreateUser(ctx, user)
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Find user by email
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err // User not found or other error
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err // Invalid password
	}

	// Generate access token (15 minutes)
	accessToken, err := s.generateToken(user.ID.Hex(), 15*time.Minute)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (7 days)
	refreshToken, err := s.generateToken(user.ID.Hex(), 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) generateToken(userID string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(expiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
