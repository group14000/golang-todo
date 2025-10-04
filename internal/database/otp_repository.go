package database

import (
	"context"
	"time"

	"github.com/group14000/golang-todo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OTPRepository interface {
	Create(ctx context.Context, otp *models.OTP) error
	FindValidOTP(ctx context.Context, email, code string, otpType models.OTPType) (*models.OTP, error)
	MarkAsUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type otpRepository struct {
	collection *mongo.Collection
}

func NewOTPRepository(client *mongo.Client) OTPRepository {
	return &otpRepository{collection: client.Database("golang-todo").Collection("otps")}
}

func (r *otpRepository) Create(ctx context.Context, otp *models.OTP) error {
	_, err := r.collection.InsertOne(ctx, otp)
	return err
}

func (r *otpRepository) FindValidOTP(ctx context.Context, email, code string, otpType models.OTPType) (*models.OTP, error) {
	var otp models.OTP
	filter := bson.M{
		"email":      email,
		"code":       code,
		"type":       otpType,
		"is_used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}
	err := r.collection.FindOne(ctx, filter).Decode(&otp)
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *otpRepository) MarkAsUsed(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"is_used": true}})
	return err
}

func (r *otpRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": time.Now()}})
	return err
}
