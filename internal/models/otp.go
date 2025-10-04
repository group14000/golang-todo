package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTPType string

const (
	OTPTypeSignup         OTPType = "signup"
	OTPTypeForgotPassword OTPType = "forgot_password"
)

type OTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Code      string             `bson:"code" json:"code"`
	Type      OTPType            `bson:"type" json:"type"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	IsUsed    bool               `bson:"is_used" json:"is_used"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
