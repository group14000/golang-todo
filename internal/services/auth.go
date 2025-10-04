package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/group14000/golang-todo/internal/database"
	"github.com/group14000/golang-todo/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     database.UserRepository
	otpRepo      database.OTPRepository
	emailService *EmailService
	jwtSecret    string
}

func NewAuthService(userRepo database.UserRepository, otpRepo database.OTPRepository, emailService *EmailService, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
	}
}

func (s *AuthService) SignUp(ctx context.Context, name, email, password string) error {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindUserByEmail(ctx, email)
	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	// Generate OTP
	otpCode := s.emailService.GenerateOTP()

	// Create OTP record
	otp := &models.OTP{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Code:      otpCode,
		Type:      models.OTPTypeSignup,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	// Store OTP
	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return err
	}

	// Send OTP email
	if err := s.emailService.SendOTP(email, otpCode, "signup"); err != nil {
		return err
	}

	return nil
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err // User not found or other error
	}

	// Check if user is verified
	if !user.IsVerified {
		return nil, fmt.Errorf("please verify your email first")
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

func (s *AuthService) VerifyOTP(ctx context.Context, email, code, name, password string) error {
	// Find valid OTP
	otp, err := s.otpRepo.FindValidOTP(ctx, email, code, models.OTPTypeSignup)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := &models.User{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Email:      email,
		Password:   string(hashedPassword),
		IsVerified: true,
		CreatedAt:  time.Now(),
	}

	// Save user to database
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	// Mark OTP as used
	return s.otpRepo.MarkAsUsed(ctx, otp.ID.Hex())
}

func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Check if user exists
	_, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Generate OTP
	otpCode := s.emailService.GenerateOTP()

	// Create OTP record
	otp := &models.OTP{
		ID:        primitive.NewObjectID(),
		Email:     email,
		Code:      otpCode,
		Type:      models.OTPTypeForgotPassword,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	// Store OTP
	if err := s.otpRepo.Create(ctx, otp); err != nil {
		return err
	}

	// Send OTP email
	return s.emailService.SendOTP(email, otpCode, "forgot_password")
}

func (s *AuthService) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	// Find valid OTP
	otp, err := s.otpRepo.FindValidOTP(ctx, email, code, models.OTPTypeForgotPassword)
	if err != nil {
		return fmt.Errorf("invalid or expired OTP")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Get user by email to get user ID
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Update password in database
	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(hashedPassword)); err != nil {
		return err
	}

	// Mark OTP as used
	return s.otpRepo.MarkAsUsed(ctx, otp.ID.Hex())
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*models.User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := s.userRepo.FindUserByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
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
