package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/group14000/golang-todo/internal/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(s *services.TodoService) *TodoHandler {
	return &TodoHandler{service: s}
}

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func (h *TodoHandler) Create(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	uid, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v := validator.New()
	if err := v.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.service.Create(c.Request.Context(), uid, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create todo"})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) List(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	uid, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	todos, err := h.service.List(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list todos"})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Get(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	uid, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	id := c.Param("id")
	tid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	todo, err := h.service.Get(c.Request.Context(), uid, tid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	uid, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	id := c.Param("id")
	tid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	var req UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title == nil && req.Description == nil && req.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	if err := h.service.Update(c.Request.Context(), uid, tid, req.Title, req.Description, req.Completed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *TodoHandler) Delete(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	uid, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	id := c.Param("id")
	tid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uid, tid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Simple health handler if needed
func (h *TodoHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
}
