package services

import (
	"context"
	"time"

	"github.com/group14000/golang-todo/internal/database"
	"github.com/group14000/golang-todo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService struct {
	repo database.TodoRepository
}

func NewTodoService(repo database.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) Create(ctx context.Context, userID primitive.ObjectID, title, description string) (*models.Todo, error) {
	todo := &models.Todo{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *TodoService) List(ctx context.Context, userID primitive.ObjectID) ([]*models.Todo, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *TodoService) Get(ctx context.Context, userID, todoID primitive.ObjectID) (*models.Todo, error) {
	return s.repo.GetByID(ctx, userID, todoID)
}

func (s *TodoService) Update(ctx context.Context, userID, todoID primitive.ObjectID, title, description *string, completed *bool) error {
	update := bson.M{"updated_at": time.Now()}
	if title != nil {
		update["title"] = *title
	}
	if description != nil {
		update["description"] = *description
	}
	if completed != nil {
		update["completed"] = *completed
	}
	return s.repo.Update(ctx, userID, todoID, update)
}

func (s *TodoService) Delete(ctx context.Context, userID, todoID primitive.ObjectID) error {
	return s.repo.Delete(ctx, userID, todoID)
}
