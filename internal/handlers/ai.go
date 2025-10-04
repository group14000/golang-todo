package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/group14000/golang-todo/internal/services"
)

type AIHandler struct {
	service *services.AIService
}

func NewAIHandler(s *services.AIService) *AIHandler { return &AIHandler{service: s} }

type AIChatRequest struct {
	Prompt   string               `json:"prompt"`
	Messages []services.AIMessage `json:"messages"` // optional multi-turn; if provided overrides Prompt
	Stream   bool                 `json:"stream"`   // request streaming
}

type AIChatResponse struct {
	Answer string `json:"answer"`
}

// @Summary      AI Chat
// @Description  Sends a prompt or multi-turn messages to the AI model. Set stream=true for SSE streaming (Swagger UI only shows non-stream examples).
// @Tags         ai
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      AIChatRequestDTO  true  "AI chat request"
// @Success      200      {object}  AIChatResponseDTO
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      502      {object}  ErrorResponse
// @Router       /ai/chat [post]
func (h *AIHandler) Chat(c *gin.Context) {
	start := time.Now()
	userID := c.GetString("user_id")
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build messages slice
	messages := req.Messages
	if len(messages) == 0 {
		if req.Prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "prompt or messages required"})
			return
		}
		messages = []services.AIMessage{{Role: "user", Content: req.Prompt}}
	}

	// Basic hash key for caching (computed here for logging; cache logic added in service later)
	cacheKeyRaw := struct{ Count int }{Count: len(messages)}
	_ = cacheKeyRaw
	// Metrics
	totalChars := 0
	for _, m := range messages {
		totalChars += len(m.Content)
	}

	if req.Stream {
		// Streaming path
		streamCh, err := h.service.ChatStream(c.Request.Context(), messages)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Stream(func(w io.Writer) bool {
			chunk, ok := <-streamCh
			if !ok {
				return false
			}
			if chunk.Err != nil {
				c.SSEvent("error", gin.H{"error": chunk.Err.Error()})
				return false
			}
			c.SSEvent("token", gin.H{"data": chunk.Text})
			return true
		})
		return
	}

	answer, err := h.service.Chat(c.Request.Context(), messages)
	latency := time.Since(start)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error(), "latency_ms": latency.Milliseconds()})
		return
	}

	// Log minimal metrics (stdout via Gin logger already). Here we set headers for observability.
	sha := sha256.Sum256([]byte(userID + time.Now().Format(time.RFC3339Nano)))
	c.Header("X-AI-Prompt-Chars", intToString(totalChars))
	c.Header("X-AI-Msg-Count", intToString(len(messages)))
	c.Header("X-AI-Latency-MS", intToString(int(latency.Milliseconds())))
	c.Header("X-AI-Req", hex.EncodeToString(sha[:6]))

	c.JSON(http.StatusOK, AIChatResponse{Answer: answer})
}

func intToString(i int) string { return strconv.FormatInt(int64(i), 10) }
