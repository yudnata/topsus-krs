# Backend API (Go Fiber v3)

A high-performance backend API built with Go Fiber v3 using a **Modular Feature-First** architecture.

## Folder Structure

```bash
backend/
├── cmd/
│   └── main.go                  # Main entry point
├── internal/
│   ├── config/                  # App configuration & env loader
│   ├── cache/                   # Redis client, keys, invalidator
│   ├── database/                # DB Connection & Auto-migrations
│   ├── modules/                 # Feature-based modules (Domains)
│   │   ├── auth/                # Auth: handler → service → repository
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   ├── routes.go
│   │   │   ├── middleware.go    # JWT guard (per modul, bukan folder global)
│   │   │   └── types.go
│   │   └── (slice 4+) admin/, mahasiswa/, pengajuan_krs/, ...
│   └── router/                  # Central router orchestration
├── pkg/                         # Shared utilities
│   ├── response/                # JSON response helper
│   └── validator/               # Input validation helper
├── .env                         # Local configuration
└── Dockerfile                   # Multi-stage build definition
```

## Architectural Workflow

This project follows a **Feature-First / Vertical Slice Architecture**:
Each feature is self-contained within its own folder under `internal/modules/`.

Data flow:
`HTTP Request` → `feature/handler` → `feature/service` → `feature/repository` → `PostgreSQL`

- **model/types**: Defines domain entities and feature-specific DTOs.
- **handler**: Manages input (parsing JSON/params) and output (status codes, standard responses).
- **service**: Contains core business logic (validation, computation, repository coordination).
- **repository**: Pure database operations using Raw SQL with `pgx`.

## Routing & Integration

Every feature has a `routes.go` file to define its internal endpoints. These modules are then registered in `internal/router/router.go`:

```go
// internal/modules/auth/routes.go
func RegisterRoutes(router fiber.Router, h *Handler, svc *Service) {
    g := router.Group("/auth")
    g.Post("/login", h.Login)
    g.Get("/profile", RequireAuth(svc), h.Profile)
}

// internal/router/router.go — hanya orchestration, tanpa business logic
func Setup(app *fiber.App, authH *auth.Handler, authSvc *auth.Service) {
    api := app.Group("/api")
    auth.RegisterRoutes(api, authH, authSvc)
}
```

## Local Setup

1. Copy `.env.example` to `.env` and configure your credentials.
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the development server:
   ```bash
   go run cmd/main.go
   ```

## API Endpoints

- **Auth**
  - `POST /api/auth/register` - Create user (role: ADMIN|MAHASISWA|DOSEN|STAFF)
  - `POST /api/auth/login` - JWT + user profile
  - `GET /api/auth/profile` - Profil user (Bearer token)
- **Health**
  - `GET /api/health` - Status API + Redis
- **Cache (uji / dev)**
  - `GET /api/cache/status` - Redis ping
  - `GET /api/cache/demo?key=` - Uji hit/miss (`X-Cache` header)
  - `POST /api/cache/demo?key=&value=` - Set cache TTL 120s
  - `DELETE /api/cache/demo?key=` - Invalidate key
