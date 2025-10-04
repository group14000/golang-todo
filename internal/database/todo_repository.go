package database

import (
	"context"

	"github.com/group14000/golang-todo/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) error
	ListByUser(ctx context.Context, userID primitive.ObjectID) ([]*models.Todo, error)
	GetByID(ctx context.Context, userID, todoID primitive.ObjectID) (*models.Todo, error)
	Update(ctx context.Context, userID, todoID primitive.ObjectID, update bson.M) error
	Delete(ctx context.Context, userID, todoID primitive.ObjectID) error
}

type todoRepository struct {
	collection *mongo.Collection
}

func NewTodoRepository(client *mongo.Client) TodoRepository {
	return &todoRepository{collection: client.Database("golang-todo").Collection("todos")}
}

func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) error {
	_, err := r.collection.InsertOne(ctx, todo)
	return err
}

func (r *todoRepository) ListByUser(ctx context.Context, userID primitive.ObjectID) ([]*models.Todo, error) {
	cur, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var todos []*models.Todo
	for cur.Next(ctx) {
		var t models.Todo
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		todos = append(todos, &t)
	}
	return todos, cur.Err()
}

func (r *todoRepository) GetByID(ctx context.Context, userID, todoID primitive.ObjectID) (*models.Todo, error) {
	var todo models.Todo
	err := r.collection.FindOne(ctx, bson.M{"_id": todoID, "user_id": userID}).Decode(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Update(ctx context.Context, userID, todoID primitive.ObjectID, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": todoID, "user_id": userID}, bson.M{"$set": update})
	return err
}

func (r *todoRepository) Delete(ctx context.Context, userID, todoID primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": todoID, "user_id": userID})
	return err
}
