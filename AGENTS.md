# reflect-backend

Go 1.25.0 REST API built with Gin, GORM (PostgreSQL), JWT auth, Cloudinary.

## Entrypoint

`cmd/api/main.go` — single-binary module, no monorepo.

## Commands

| Action | Command |
|---|---|
| Run | `go run cmd/api/main.go` |
| Build | `go build -o main cmd/api/main.go` |
| Deps | `go mod tidy` |
| Hot reload | `air` (requires `.air.toml` — create from template; see `.gitignore`) |

No tests exist yet. No CI/CD, no Docker, no Makefile.

## Architecture

- `internal/modules/<name>/` — each module: `model`, `repository`, `service`, `handler`, (request/response DTOs)
- Wiring happens in `internal/routes/routes.go` — construct repository → service → handler
- `internal/database/postgres.go` — global `DB *gorm.DB`, auto-migrates all models on startup
- `internal/config/config.go` — loads via `godotenv`, sensible defaults
- `internal/utils/errors.go` — `AppError` with typed constructors (`BadRequest`, `Unauthorized`, etc.)
- `internal/middleware/auth_middleware.go` — JWT parse + `RoleMiddleware(allowedRoles...)`

## API

All routes under `/api/v1`:
- **Public** — `auth/` (register/login), `products/`, `categories/`, `collections/`
- **Protected** (JWT required) — `/me`, `/cart`, `/wishlist`, `/addresses`, `/orders` (own), `/checkout`, `/upload`
- **Admin** (JWT + role=admin) — `/admin/products`, `/admin/categories`, `/admin/collections`, `/admin/orders/:id/status`

## Key conventions

- UUID PKs with `BeforeCreate` hook generating new UUIDs
- Roles: `"user"` (default), `"admin"`
- `middleware.GetCurrentUserID(c)` to get authenticated user ID in handlers
- Handlers call `c.Error(err)` on service errors; `ErrorMiddleware` converts `*AppError` to JSON
- Validation errors use `utils.ValidationErrorResponse(c, err.Error())` directly (not `c.Error`)
- Uploads use Cloudinary (configured via env vars)

## Setup

Copy `.env.example` → `.env` and fill in:
- PostgreSQL connection params (default: localhost:5432, `reflect_db`)
- Cloudinary credentials (for image upload)
- `JWT_SECRET` (default: `"secret"`)

Database auto-migrates on startup — no manual migration step.
