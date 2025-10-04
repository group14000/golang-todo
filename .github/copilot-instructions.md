# Golang Todo API - AI Agent Instructions

## Architecture Overview

This is a **Clean Architecture** Go REST API with **OTP-based authentication** and **JWT authorization**.

### Core Layers
- **`cmd/server/`**: Application entry point with dependency injection
- **`internal/models/`**: Domain entities (`User`, `Todo`, `OTP`) with MongoDB ObjectID types
- **`internal/database/`**: Repository pattern with MongoDB queries scoped by `user_id`
- **`internal/services/`**: Business logic layer handling auth flows and CRUD operations
- **`internal/handlers/`**: HTTP controllers using Gin framework
- **`internal/middleware/`**: JWT bearer token validation
- **`api/`**: Route definitions separating public/protected endpoints

### Authentication Flow (Critical)
1. **Signup**: `POST /signup` → sends OTP via email (user NOT created yet)
2. **Verify**: `POST /verify-otp` → validates OTP + creates verified user account
3. **Login**: `POST /login` → requires `IsVerified=true`, returns JWT tokens
4. **Reset**: `POST /forgot-password` → `POST /reset-password` (OTP-based)

**Key Security Pattern**: All todos are scoped by `userID` from JWT claims in middleware.

## Development Workflows

### Build & Run
```bash
go run ./cmd/server          # Development server on :8080
go build ./cmd/server        # Compile binary
```

### Environment Setup
- Copy `.env` file with MongoDB URL, JWT secret, and Gmail SMTP config
- MongoDB collections: `users`, `todos`, `otps` (auto-created)
- OTPs expire in 10 minutes and are single-use

## Project-Specific Patterns

### Repository Pattern
```go
// All repos follow this interface pattern
type TodoRepository interface {
    GetByID(ctx, userID, todoID primitive.ObjectID) (*models.Todo, error)
}
// Note: Always pass userID for data isolation
```

### Service Injection Chain
```go
// Services depend on multiple repositories + external services
authService := services.NewAuthService(userRepo, otpRepo, emailService, jwtSecret)
```

### JWT Middleware Pattern
```go
// Middleware extracts user_id and sets in gin.Context
userID := c.GetString("user_id")  // Available in all protected handlers
```

### Error Handling Convention
- Services return `fmt.Errorf()` with business context
- Handlers map to appropriate HTTP status codes
- MongoDB errors bubble up (no custom wrapping)

## Key Integration Points

### MongoDB Dependencies
- All models use `primitive.ObjectID` for IDs
- BSON tags for field mapping: `bson:"user_id" json:"user_id"`
- Time fields auto-set in services (`CreatedAt`, `UpdatedAt`, `ExpiresAt`)

### Email Service (External)
- Uses `gopkg.in/gomail.v2` with Gmail SMTP
- HTML templates embedded in service methods
- 6-digit crypto/rand OTP generation

### Validation Pattern
- Struct tags: `validate:"required,email,min=6,len=6"`
- Per-handler validator instances (not global)

## Common Gotchas

1. **ObjectID Conversion**: Always validate `primitive.ObjectIDFromHex()` in handlers
2. **User Verification**: Login fails if `IsVerified=false` (check signup flow)
3. **CORS/Auth Headers**: Protected routes require `Authorization: Bearer <token>`
4. **OTP Lifecycle**: OTPs are marked `is_used=true` after verification
5. **Password Security**: Never return password in JSON (`json:"-"` tag)

## Testing Approach
- No unit tests currently (add to `*_test.go` files)
- Manual API testing via Postman/curl
- MongoDB data verification for complex flows