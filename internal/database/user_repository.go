package database

import (
	"context"

	"github.com/group14000/golang-todo/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) UserRepository {
	return &userRepository{
		collection: client.Database("golang-todo").Collection("users"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
