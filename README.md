# Golang Todo API

A clean-architecture Go REST API featuring OTP-based signup verification, JWT auth, Todo CRUD, AI chat (OpenRouter DeepSeek), and generated Swagger documentation.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25-blue" />
  <img src="https://img.shields.io/badge/Framework-Gin-green" />
  <img src="https://img.shields.io/badge/DB-MongoDB-success" />
  <img src="https://img.shields.io/badge/Auth-JWT-orange" />
  <img src="https://img.shields.io/badge/AI-OpenRouter-purple" />
  <img src="https://img.shields.io/badge/Docs-Swagger-lightgrey" />
</p>

## ✨ Features
- OTP email verification flow (deferred user creation)
- Secure JWT (access + refresh) authentication
- Password reset via OTP
- User profile endpoint
- Todo CRUD scoped per-user (Mongo isolation)
- AI chat endpoint (multi-turn + optional streaming via SSE)
- Structured validation & consistent error schema
- Auto-generated Swagger docs (`/swagger/index.html`)
- Clean layering (Handlers → Services → Repositories → Mongo)

## 🗂 Project Structure
```
cmd/server/             # Composition root (wiring, swagger serve)
api/routes.go           # Public vs protected route groups
internal/models/        # Domain models (User, Todo, OTP)
internal/database/      # Mongo repositories
internal/services/      # Business logic (Auth, Email, Todo, AI)
internal/handlers/      # Gin handlers + DTOs + swagger annotations
internal/middleware/    # JWT auth middleware
.docs/ / docs/          # (Generated) swagger spec (do not hand edit)
.github/copilot-instructions.md  # Agent guide
```

## 🔐 Auth & OTP Flow
1. `POST /signup` — store OTP, email it (no user yet)
2. `POST /verify-otp` — validate OTP, create verified user
3. `POST /login` — return access + refresh tokens
4. `POST /forgot-password` — issue reset OTP
5. `POST /reset-password` — validate OTP & update password
6. `GET /profile` — return current user (requires Bearer token)

> All protected endpoints require: `Authorization: Bearer <access_token>`

## 🤖 AI Chat
Endpoint: `POST /ai/chat`
Body options:
```json
{
  "prompt": "Explain clean architecture in Go"
}
```
Multi-turn & streaming:
```json
{
  "messages": [
    {"role": "user", "content": "You are a helpful assistant"},
    {"role": "user", "content": "Summarize hexagonal architecture"}
  ],
  "stream": true
}
```
Streaming returns Server-Sent Events (`token` events). Swagger UI shows only non-stream version.

Response (non-stream):
```json
{ "answer": "Clean Architecture in Go involves ..." }
```

### AI Headers
| Header | Meaning |
|--------|---------|
| X-AI-Prompt-Chars | Total input characters |
| X-AI-Msg-Count | Number of messages sent |
| X-AI-Latency-MS | Round-trip latency |
| X-AI-Req | Short request fingerprint |

## 📦 Environment Variables
Create a `.env` file:
```
MONGODB_URI=mongodb://localhost:27017/todoapp
JWT_SECRET=change-me
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your@gmail.com
SMTP_PASS=app-password
AI_API_KEY=sk-or-openrouter-key
```

## 🚀 Run
```bash
go run ./cmd/server        # Dev run (http://localhost:8080)
```
Swagger: http://localhost:8080/swagger/index.html

Build binary:
```bash
go build ./cmd/server
```

## 🧪 Quick cURL Examples
Signup & Verify:
```bash
curl -X POST http://localhost:8080/signup -H "Content-Type: application/json" -d '{"name":"Jane","email":"jane@example.com","password":"Secretp@ss1"}'
# Check email for OTP
curl -X POST http://localhost:8080/verify-otp -H "Content-Type: application/json" -d '{"name":"Jane","email":"jane@example.com","password":"Secretp@ss1","otp":"123456"}'
```
Login:
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/login -H "Content-Type: application/json" -d '{"email":"jane@example.com","password":"Secretp@ss1"}' | jq -r .access_token)
```
Create Todo:
```bash
curl -X POST http://localhost:8080/todos -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"title":"Buy milk","description":"2L"}'
```
AI Chat:
```bash
curl -X POST http://localhost:8080/ai/chat -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"prompt":"Explain Go contexts"}'
```
Streaming AI Chat (SSE):
```bash
curl -N -X POST http://localhost:8080/ai/chat -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"prompt":"Stream tokens","stream":true}'
```

## 🧱 Design Principles
- **Isolation by user**: every todo query includes `userID` filter
- **Deferred user creation**: only after OTP verification
- **Stateless API**: tokens contain all auth context
- **Separation of concerns**: handlers only orchestrate & validate
- **Regenerate docs**: run `swag init -g cmd/server/main.go -o docs` after changing annotations

## 🛠 Extending
| Add | Steps |
|-----|-------|
| New endpoint | DTO → handler + swagger comments → service → repo (if needed) → route → `swag init` |
| New model | Add struct in `internal/models` + repository + service logic |
| AI rate limiting | Implement token bucket in services or middleware; expose 429 on exceed |
| AI caching | Hash(model+messages) → store answer (skip when `stream=true`) |

## 🧭 Troubleshooting
| Problem | Fix |
|---------|-----|
| Empty Swagger paths | Annotations inside function body → move ABOVE func |
| 500 on ObjectID | Add `primitive.ObjectIDFromHex` validation in handler |
| AI error 502 | Missing/invalid `AI_API_KEY` or upstream failure |
| OTP always invalid | Expired (>10m) or already `IsUsed=true` |

## 📌 Roadmap
- Per‑user AI rate limiting
- AI response caching (5m TTL)
- Refresh token rotation endpoint
- Unit/integration tests (services & repositories)

## 📄 License
MIT (add LICENSE file if distributing publicly)

---
**Happy building!** Contributions welcome—follow established patterns and keep layers clean.
