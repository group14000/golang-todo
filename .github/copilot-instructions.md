type TodoRepository interface {
authService := services.NewAuthService(userRepo, otpRepo, emailService, jwtSecret)
# Golang Todo API – AI Agent Guide

Concise operational knowledge for this repository. Follow existing patterns; do not introduce new frameworks without need.

## 1. Architecture & Layering
- `cmd/server/`: Composition root (builds repositories, services, handlers, middleware, swagger binding).
- `internal/models/`: Plain structs with BSON/JSON tags. Never leak passwords (User.Password has `json:"-"`).
- `internal/database/`: Mongo repositories. Every data access for user‑scoped entities requires the caller’s `userID` (ObjectID) to enforce isolation.
- `internal/services/`: Business logic (Auth, OTP, Email, Todo, AI). Services return Go errors (no custom error types) and do not write HTTP responses.
- `internal/handlers/`: Gin HTTP layer: validate input, map service errors → status codes, perform ObjectID parsing.
- `internal/middleware/`: JWT middleware extracts `user_id` claim and sets it in `gin.Context`.
- `api/routes.go`: Central route wiring (public vs protected groups).
- `docs/`: Generated Swagger (`swag init`). Do not hand edit except for quick host tweaks.

## 2. Auth & OTP Flow (Critical)
1. `POST /signup`: store OTP only (no user yet) → Email via `EmailService`.
2. `POST /verify-otp`: validates OTP, creates verified user (IsVerified=true), marks OTP used.
3. `POST /login`: requires verified user; returns Access + Refresh (both JWT HS256; refresh not yet rotated anywhere else).
4. Password reset: `/forgot-password` issues OTP type `forgot_password`, `/reset-password` validates & updates hash.

Security rule: All Todo & AI operations rely on `user_id` from JWT middleware; repositories always take `userID` to prevent cross-user access.

## 3. Data & Models
- IDs: `primitive.ObjectID`. Always validate with `primitive.ObjectIDFromHex` before passing to repos.
- Time fields (CreatedAt, UpdatedAt, ExpiresAt) set in services; keep handlers slim.
- OTP: single-use, 10‑minute expiry, marked used immediately after success.

## 4. AI Chat Service
- File: `internal/services/ai.go` provides `Chat` (non-stream) and `ChatStream` (SSE) against OpenRouter model `deepseek/deepseek-chat-v3.1:free` (fixed).
- Handler accepts either `prompt` or `messages[]` plus `stream` flag. Multi-turn supported by passing array of `{role, content}`.
- Streaming uses `text/event-stream` and emits `token` events; Swagger documents only non-stream use.
- Observability headers set: `X-AI-Prompt-Chars`, `X-AI-Msg-Count`, `X-AI-Latency-MS`, `X-AI-Req`.
- Pending (not yet implemented): per-user rate limiting & in-memory cache (TTL) — add in services layer or dedicated package; expose metrics sparingly.

## 5. Swagger / API Docs
- Generated via `swag init -g cmd/server/main.go -o docs`.
- Annotations live immediately ABOVE handler functions (not inside bodies) or they are ignored.
- Protected routes declare `@Security BearerAuth`. JWT header format: `Authorization: Bearer <token>`.
- Do not edit `docs/docs.go` manually (regeneration overwrites); if build fails due to `LeftDelim/RightDelim`, remove those fields or upgrade swag.

## 6. Validation & Error Conventions
- Use `validator.New()` inside each handler (no shared global). Required & email patterns via struct tags.
- On invalid input: return `400` with `{ "error": message }`.
- Auth failures: `401`; missing/invalid resource: `404`; service/internal issues: `500` (some existing handlers use `502` for AI upstream failures).
- Keep error payload shape consistent using `ErrorResponse` DTO (defined in `handlers/dto.go`).

## 7. Email & OTP
- `EmailService` uses Gomail + SMTP env vars. OTP generated via crypto/rand style numeric string (6 digits). If you add new OTP types, extend enum in model & service switch logic.

## 8. Dependency Injection Pattern
Example (`cmd/server/main.go`):
```go
userRepo := database.NewUserRepository(client)
otpRepo  := database.NewOTPRepository(client)
emailSvc := services.NewEmailService(cfg)
authSvc  := services.NewAuthService(userRepo, otpRepo, emailSvc, cfg.JWTSecret)
todoSvc  := services.NewTodoService(todoRepo)
aiSvc    := services.NewAIService(cfg.AIAPIKey)
```
Add new services by extending this wiring; avoid global singletons.

## 9. Build & Run
```bash
go run ./cmd/server    # starts on :8080
go build ./cmd/server  # builds binary
```
Swagger UI: `http://localhost:8080/swagger/index.html` (doc JSON at `/swagger/doc.json`).

## 10. Adding New Endpoints (Checklist)
1. Define request/response DTO (exported if you want Swagger schema).
2. Add handler method (validate → service → map error → response).
3. Add route in `api/routes.go` (place under protected group if token-required).
4. Add Swagger annotations ABOVE the function.
5. Regenerate docs with `swag init`.

## 11. Common Pitfalls
- Placing Swagger annotations inside function bodies → missing paths.
- Forgetting ObjectID validation leads to 500 from Mongo; always pre-check and return 400.
- Returning password or internal fields: ensure struct JSON tags enforce omissions.
- Missing AI_API_KEY env disables AI; `Chat` returns config error.
- Deleting or editing `docs/docs.go` manually without regen causes drift.

## 12. Future Enhancements (Referenced but Pending)
- AI rate limiting (token bucket per user, likely in-memory map with mutex).
- AI response cache (hash(model+messages) → answer, 5m TTL, skip for streaming).
- Refresh token rotation endpoint (currently not implemented).

Keep instructions concise—opt for small service-focused changes, preserve existing layering, and never bypass repositories for Mongo access from handlers.