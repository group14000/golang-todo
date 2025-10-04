package services

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"

	"github.com/group14000/golang-todo/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

func (s *EmailService) SendOTP(to, otp string, otpType string) error {
	subject := "Verify Your Email"
	body := fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Your OTP for %s is: <strong>%s</strong></p>
		<p>This code will expire in 10 minutes.</p>
		<p>If you didn't request this, please ignore this email.</p>
	`, otpType, otp)

	if otpType == "forgot_password" {
		subject = "Reset Your Password"
		body = fmt.Sprintf(`
			<h2>Password Reset</h2>
			<p>Your OTP for password reset is: <strong>%s</strong></p>
			<p>This code will expire in 10 minutes.</p>
			<p>If you didn't request this, please ignore this email.</p>
		`, otp)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.EmailUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.EmailHost, s.config.EmailPort, s.config.EmailUser, s.config.EmailPassword)
	if s.config.EmailUseTLS {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return d.DialAndSend(m)
}

func (s *EmailService) GenerateOTP() string {
	const digits = "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		otp[i] = digits[num.Int64()]
	}
	return string(otp)
}
