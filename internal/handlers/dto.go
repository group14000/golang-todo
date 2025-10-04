package handlers

// ErrorResponse represents a standard error payload
// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

// SignupRequestDTO represents signup request
// swagger:model SignupRequest
type SignupRequestDTO struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"Secretp@ss1"`
}

// VerifyOTPRequestDTO represents verify-otp request
// swagger:model VerifyOTPRequest
type VerifyOTPRequestDTO struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"Secretp@ss1"`
	OTP      string `json:"otp" example:"123456"`
}

// LoginRequestDTO represents login request
// swagger:model LoginRequest
type LoginRequestDTO struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"Secretp@ss1"`
}

// ForgotPasswordRequestDTO represents forgot password request
// swagger:model ForgotPasswordRequest
type ForgotPasswordRequestDTO struct {
	Email string `json:"email" example:"john@example.com"`
}

// ResetPasswordRequestDTO represents reset password request
// swagger:model ResetPasswordRequest
type ResetPasswordRequestDTO struct {
	Email       string `json:"email" example:"john@example.com"`
	OTP         string `json:"otp" example:"123456"`
	NewPassword string `json:"new_password" example:"NewSecretp@ss1"`
}

// CreateTodoRequestDTO represents create todo request
// swagger:model CreateTodoRequest
type CreateTodoRequestDTO struct {
	Title       string `json:"title" example:"Buy milk"`
	Description string `json:"description" example:"2 liters of whole milk"`
}

// UpdateTodoRequestDTO represents update todo request
// swagger:model UpdateTodoRequest
type UpdateTodoRequestDTO struct {
	Title       *string `json:"title" example:"Buy bread"`
	Description *string `json:"description" example:"Whole grain"`
	Completed   *bool   `json:"completed" example:"true"`
}

// AIChatMessageDTO represents a single AI chat message
// swagger:model AIChatMessage
type AIChatMessageDTO struct {
	Role    string `json:"role" example:"user"`
	Content string `json:"content" example:"Hello, how are you?"`
}

// AIChatRequestDTO represents AI chat request
// swagger:model AIChatRequest
type AIChatRequestDTO struct {
	Prompt   string             `json:"prompt" example:"Explain clean architecture in Go"`
	Messages []AIChatMessageDTO `json:"messages"`
	Stream   bool               `json:"stream" example:"false"`
}

// AIChatResponseDTO represents AI chat response
// swagger:model AIChatResponse
type AIChatResponseDTO struct {
	Answer string `json:"answer" example:"Clean Architecture in Go involves..."`
}
