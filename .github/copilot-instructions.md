# Reborn Project - AI Coding Assistant Instructions

## Architecture Overview

This is a **Go + Echo backend with React frontend** monorepo serving as a web platform with user management. Key architectural patterns:

- **Service-oriented backend**: Uses `ServiceManager` pattern with dependency injection
- **gRPC microservice client**: Auth service communicates via gRPC to external user service
- **Single-page app routing**: Frontend handles routing, backend serves static files + API
- **Configuration-driven**: Uses TOML config with Viper for environment-specific settings

## Project Structure & Key Files

```
cmd/main.go              # Application entry point with graceful shutdown
internal/
  services/              # Business logic layer (AuthService, ServiceManager)  
  routers/              # Route definitions (api_v1.go, auth.go, page.go)
  handlers/             # HTTP request handlers
  middlewares/          # Echo middleware (auth, CORS, rate limiting)
  client/               # gRPC client for external services
configs/                # Configuration management
website/                # React SPA (TypeScript + Vite + Tailwind)
  src/api/              # Auto-generated TypeScript API client
```

## Development Workflows

### API Development
```bash
# Generate Swagger docs + TypeScript client (REQUIRED after API changes)
make swag

# Build for development (includes swag generation)
make build-dev

# Build frontend only
make website
```

### Code Quality & Formatting
```bash
# Format Go code (golines + gofumpt)
make fmt

# Lint and auto-fix Go code (golangci-lint)
make lint

# Install development tools
make install
```

### Key Patterns

**Service Manager Pattern**: All services initialized through `ServiceManager` - always use dependency injection, never direct instantiation.

**Route Registration**: Three route groups in `main.go`:
- `RegisterAPIv1Routes` - REST API with auth middleware
- `RegisterAuthRoutes` - OAuth callback endpoints  
- `RegisterPageRoutes` - Serves React SPA + static files

**Authentication Flow**:
- Session-based auth via `middlewares.LoginSession()`
- Admin routes require `middlewares.AdminOnly()` 
- Frontend handles OAuth redirects, backend validates sessions

**Configuration**: Use `config.Load()` pattern, never hardcode values. All config keys defined as constants in `configs/config.go`.

**Logging**: Use structured logging with `slog.ErrorContext()`, `slog.InfoContext()`, etc. Always pass request context for tracing:
```go
slog.ErrorContext(c.Request().Context(), "Failed to close file", 
    "file", filePath, 
    "error", closeErr)
```
Never use Echo's `c.Logger()` - use Go's standard `log/slog` package instead.

## API Generation Pipeline

Critical: API changes require **regeneration**:
1. Update Go handlers with Swagger annotations
2. Run `make swag` - generates OpenAPI spec + TypeScript client
3. Frontend imports from `@/api/api` (auto-generated)

## Frontend-Backend Integration

- **Static serving**: `/admin/*` routes serve React app, client-side routing handles navigation
- **API client**: Auto-generated from OpenAPI spec, uses Axios with proper typing
- **Auth context**: React `useAuth` hook syncs with backend session state
- **Admin protection**: Both backend middleware and frontend route guards

## Docker & Deployment

Multi-stage build:
1. Node.js stage builds React frontend
2. Go stage builds backend binary  
3. Ubuntu runtime stage combines both

Configuration via `configs/default.toml` - modify for different environments.

## Common Gotchas

- **Always run `make swag` after API changes** - TypeScript client won't reflect backend changes otherwise
- **Run `make fmt` before committing** - Uses golines + gofumpt for consistent Go formatting
- **Use `make lint` for code quality** - golangci-lint with auto-fix enabled
- **Use structured logging with context** - Always use `slog.ErrorContext(c.Request().Context(), ...)` instead of Echo's logger
- **Use ServiceManager for all service access** - don't instantiate services directly
- **Frontend routing**: Admin routes must be under `/admin/*` for backend auth middleware to work
- **CORS**: Configured in `middlewares/common.go` - update origins for production
- **Rate limiting**: Default 20 req/sec per IP - adjust in middleware if needed

## External Dependencies

- **oj-lab/user-service**: gRPC service for user management (configured via `auth_service.address`)
- **oj-lab/go-webmods**: Internal library for app initialization and config management
