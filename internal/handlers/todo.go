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

// @Summary      Create todo
// @Description  Creates a new todo item for the authenticated user.
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      CreateTodoRequestDTO  true  "Create todo"
// @Success      201      {object}  models.Todo
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /todos [post]
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

// @Summary      List todos
// @Description  Lists all todos for the authenticated user.
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Todo
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /todos [get]
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

// @Summary      Get todo
// @Description  Retrieves a single todo by ID.
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Todo ID"
// @Success      200  {object}  models.Todo
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /todos/{id} [get]
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

// @Summary      Update todo
// @Description  Partially updates a todo.
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                 true  "Todo ID"
// @Param        payload  body      UpdateTodoRequestDTO   true  "Update todo"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /todos/{id} [patch]
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

// @Summary      Delete todo
// @Description  Deletes a todo by ID.
// @Tags         todos
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Todo ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /todos/{id} [delete]
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
