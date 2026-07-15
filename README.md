# 🏋️ Habit Tracker API

A production-ready REST API for tracking daily habits, built with Go and PostgreSQL. Features JWT authentication, refresh token rotation, automated database migrations, Docker support, and a clean layered architecture following the repository pattern.

---

## ✨ Features

- **User Authentication** — Register, login, logout with secure password hashing
- **JWT Access Tokens** — Short-lived tokens for API authorization
- **Refresh Token Rotation** — Long-lived refresh tokens with revocation support
- **Habit Management** — Full CRUD for personal habits
- **Daily Check-ins** — Track habit completions with one-per-day enforcement
- **Streak Tracking** — Calculate current consecutive completion streaks
- **Dashboard** — Aggregated daily progress at a glance
- **Automated Migrations** — Schema migrations run automatically on startup
- **Docker Support** — Single-command deployment with Docker Compose
- **Graceful Shutdown** — Clean server shutdown on SIGINT / SIGTERM
- **Input Validation** — Request validation using go-playground/validator
- **Layered Architecture** — Handler → Service → Repository separation
- **Production-ready Configuration** — Environment-based configuration via `.env`

---

## 🛠 Tech Stack

| Technology | Purpose |
|---|---|
| [Go](https://go.dev/) | Backend language |
| [PostgreSQL 17](https://www.postgresql.org/) | Relational database |
| [pgx/v5](https://github.com/jackc/pgx) | PostgreSQL driver and connection pool |
| [golang-migrate](https://github.com/golang-migrate/migrate) | Database schema migrations |
| [Docker](https://www.docker.com/) | Containerization |
| [JWT](https://github.com/golang-jwt/jwt) | Access token authentication |
| [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) | Password hashing |
| [validator](https://github.com/go-playground/validator) | Struct and input validation |
| [godotenv](https://github.com/joho/godotenv) | Environment variable loading |

---

## 📁 Project Structure

```
habit-tracker/
├── cmd/
│   └── api/
│       └── main.go              # Application entrypoint
├── internal/
│   ├── app/                     # Application container
│   ├── config/                  # Environment configuration
│   ├── constants/               # Application-wide constants
│   ├── db/                      # Database connection and migrations
│   ├── dto/                     # Data Transfer Objects (request/response)
│   ├── errors/                  # Custom application errors
│   ├── handlers/                # HTTP handlers (controllers)
│   ├── middleware/              # Authentication middleware
│   ├── models/                  # Domain models
│   ├── repositories/            # Data access layer (raw SQL)
│   ├── response/                # Standardized API response helpers
│   ├── routes/                  # Route registration
│   ├── services/                # Business logic layer
│   ├── utils/                   # JWT manager and utilities
│   └── validator/               # Request validation
├── migrations/                  # SQL migration files
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # Docker Compose orchestration
├── .env.example                 # Example environment variables
└── go.mod                       # Go module definition
```

---

## 📋 Prerequisites

- [Go](https://go.dev/dl/) 1.25 or later
- [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) (for containerized setup)
- [PostgreSQL 17](https://www.postgresql.org/download/) (only if running without Docker)

---

## 🔧 Environment Variables

Create a `.env` file in the project root. See [`.env.example`](.env.example) for reference.

| Variable | Description | Example |
|---|---|---|
| `APP_PORT` | Port the API server listens on | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `your_password` |
| `DB_NAME` | Database name | `habit_tracker` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `JWT_SECRET` | Secret key for signing JWTs | `your-secret-key` |
| `ACCESS_TOKEN_TTL` | Access token time-to-live | `15m` |
| `REFRESH_TOKEN_TTL` | Refresh token time-to-live | `720h` |

> ⚠️ **Never commit your `.env` file.** It is excluded via `.gitignore`.

---

## 🚀 Running without Docker

**1. Install dependencies:**

```bash
go mod tidy
```

**2. Create the PostgreSQL database:**

```sql
CREATE DATABASE habit_tracker;
```

**3. Configure environment variables:**

```bash
cp .env.example .env
# Edit .env with your database credentials and secrets
```

**4. Start the server:**

```bash
go run ./cmd/api
```

The API will start at `http://localhost:8080`. Migrations run automatically on startup.

---

## 🐳 Running with Docker

**Start everything (PostgreSQL + API):**

```bash
docker compose up --build
```

**Run in background:**

```bash
docker compose up --build -d
```

**Stop all containers:**

```bash
docker compose down
```

**Stop and remove all data (including database volumes):**

```bash
docker compose down -v
```

> Docker Compose automatically starts PostgreSQL, waits for it to be ready, runs pending migrations, and starts the API. No manual migration commands are required.

---

## 🗄️ Database Migrations

Migrations are managed by [golang-migrate](https://github.com/golang-migrate/migrate) and run **automatically on application startup**.

- Migration files live in the root [`migrations/`](migrations/) directory.
- Each migration has an `up` (apply) and `down` (rollback) SQL file.
- On startup, the application runs all pending migrations before starting the HTTP server.
- If no pending migrations exist, `migrate.ErrNoChange` is gracefully ignored.
- Migrations are idempotent and safe for production — restarting the application never re-runs completed migrations.

**Migration files:**

```
migrations/
├── 000001_create_users.up.sql
├── 000001_create_users.down.sql
├── 000002_create_refresh_tokens.up.sql
├── 000002_create_refresh_tokens.down.sql
├── 000003_create_habits.up.sql
├── 000003_create_habits.down.sql
├── 000004_create_habit_logs.up.sql
└── 000004_create_habit_logs.down.sql
```

---

## 📡 API Endpoints

### Health Check

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Health check |

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive tokens |
| `POST` | `/api/v1/auth/refresh` | Refresh access token |
| `POST` | `/api/v1/auth/logout` | Logout (revoke refresh token) |
| `POST` | `/api/v1/auth/logout-all` | Logout from all devices |

### User

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/me` | Get authenticated user profile |

### Habits

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/habits` | Create a new habit |
| `GET` | `/api/v1/habits` | List all habits |
| `GET` | `/api/v1/habits/{id}` | Get a habit by ID |
| `PUT` | `/api/v1/habits/{id}` | Update a habit |
| `DELETE` | `/api/v1/habits/{id}` | Soft delete a habit |
| `POST` | `/api/v1/habits/{id}/check-in` | Check in a habit for today |
| `GET` | `/api/v1/habits/{id}/history` | Get habit completion history |
| `GET` | `/api/v1/habits/{id}/streak` | Get current streak |
| `GET` | `/api/v1/habits/today` | Get all habits checked in today |

### Dashboard

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/dashboard` | Get daily progress summary |

---

## 🏗 Architecture

The application follows a **layered architecture** with strict dependency flow:

```
HTTP Request
    │
    ▼
┌──────────┐
│ Handler  │  ← Decodes request, validates input, returns HTTP response
└────┬─────┘
     │
     ▼
┌──────────┐
│ Service  │  ← Contains business logic, enforces ownership rules
└────┬─────┘
     │
     ▼
┌──────────────┐
│ Repository   │  ← Executes raw SQL queries via pgx
└────┬─────────┘
     │
     ▼
┌──────────────┐
│ PostgreSQL   │
└──────────────┘
```

**Key principles:**

- **Handlers** never access the database directly.
- **Repositories** never know about HTTP concepts.
- **Services** contain all business logic and validation.
- **Context** (`context.Context`) is propagated through all layers.
- Repositories return **domain models**; Handlers convert them to **DTOs**.
- User ID is always extracted from JWT claims — never accepted from request bodies.

---

## 🔒 Security

| Measure | Details |
|---|---|
| **Password Hashing** | Passwords are hashed using `bcrypt` before storage. |
| **JWT Authentication** | Short-lived access tokens with configurable TTL. |
| **Refresh Tokens** | Hashed with SHA-256 before storage. Supports rotation and revocation. |
| **SQL Injection Protection** | All queries use parameterized statements via `pgx`. |
| **Environment Variables** | Secrets are loaded from `.env` files, never hardcoded. |
| **Ownership Enforcement** | Users can only access and modify their own resources. |

---

## 🐳 Docker

The Docker setup consists of two services:

| Service | Image | Purpose |
|---|---|---|
| `postgres` | `postgres:17` | PostgreSQL database with persistent volume |
| `api` | Custom build | Go API with embedded migrations |

- Both services share a dedicated Docker network (`habit-network`).
- PostgreSQL data is persisted using a named volume (`postgres_data`).
- The API container automatically runs migrations on startup.
- `docker compose up` is the only command needed to start the full stack.

---

## ⏹ Graceful Shutdown

The application handles `SIGINT` (Ctrl+C) and `SIGTERM` signals gracefully:

1. The server stops accepting new connections.
2. In-flight requests are given up to **15 seconds** to complete.
3. The database connection pool is closed cleanly.
4. The process exits with a success status.

This ensures zero data loss during deployments and container restarts.

---

## 🔮 Future Improvements

- [ ] Redis caching for frequently accessed data
- [ ] Email verification on registration
- [ ] Password reset flow
- [ ] Swagger / OpenAPI documentation
- [ ] Rate limiting middleware
- [ ] Structured logging (slog / zerolog)
- [ ] Prometheus metrics
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Comprehensive test suite (unit + integration)
- [ ] Habit categories and tags
- [ ] Weekly / monthly analytics

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).
