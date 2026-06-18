# Habit Tracker API

A production-style Habit Tracker backend built with **Go**, **net/http**, and **PostgreSQL**.

The goal of this project is to learn Go and its ecosystem by building a real-world backend without relying on frameworks like Gin, Echo, or Fiber.

---

## Features

### Authentication

* User registration
* User login
* JWT-based authentication
* Refresh token support
* Logout current device
* Logout from all devices
* Token rotation

### Habit Management

* Create habits
* List habits
* Update habits
* Delete habits

### Daily Check-ins

* Mark a habit as completed for the day
* Prevent duplicate check-ins

### Streak Tracking

* Current streak calculation
* Longest streak calculation

### Middleware

* JWT authentication
* Request logging
* Panic recovery
* CORS support

### Production Features

* Graceful shutdown
* Configuration using environment variables
* PostgreSQL connection pooling
* Standard project structure
* Dependency injection
* Validation
* Error handling
* Response helpers

---

# Tech Stack

* Go
* net/http
* PostgreSQL
* pgx/v5
* JWT
* bcrypt

---

# Project Structure

```text
habit-tracker/

cmd/
└── api/
    └── main.go

internal/
├── app/
├── config/
├── db/
├── handlers/
├── routes/

.env.example

go.mod
go.sum
```

---

# Getting Started

## Clone the repository

```bash
git clone https://github.com/Hettank/habit-tracker.git
cd habit-tracker
```

---

## Install dependencies

```bash
go mod tidy
```

---

## Create environment file

Copy:

```bash
cp .env.example .env
```

Fill the values:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=habit_tracker
DB_SSLMODE=disable
```

---

## Create PostgreSQL database

Login to PostgreSQL:

```bash
psql -U postgres -h localhost
```

Create database:

```sql
CREATE DATABASE habit_tracker;
```

---

## Run the application

From the project root:

```bash
go run ./cmd/api
```

---

## Health Check

Endpoint:

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

---

# Planned API Endpoints

## Auth

```http
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
POST /api/v1/auth/logout-all
```

## Habits

```http
POST   /api/v1/habits
GET    /api/v1/habits
PATCH  /api/v1/habits/{id}
DELETE /api/v1/habits/{id}
```

## Check-ins

```http
POST /api/v1/habits/{id}/checkin
```

---

# Database Schema

### users

```sql
id BIGSERIAL PRIMARY KEY
email VARCHAR(255) UNIQUE NOT NULL
password_hash TEXT NOT NULL
created_at TIMESTAMP DEFAULT NOW()
```

### habits

```sql
id BIGSERIAL PRIMARY KEY
user_id BIGINT REFERENCES users(id)
title VARCHAR(255)
current_streak INT
longest_streak INT
created_at TIMESTAMP DEFAULT NOW()
```

### habit_logs

```sql
id BIGSERIAL PRIMARY KEY
habit_id BIGINT REFERENCES habits(id)
logged_date DATE
UNIQUE(habit_id, logged_date)
```

### refresh_tokens

```sql
id BIGSERIAL PRIMARY KEY
user_id BIGINT REFERENCES users(id)
token_hash TEXT
expires_at TIMESTAMP
revoked BOOLEAN DEFAULT FALSE
created_at TIMESTAMP DEFAULT NOW()
```

---

# Development Roadmap

* [x] Project initialization
* [x] PostgreSQL connection
* [x] Graceful shutdown
* [x] Health endpoint
* [ ] Response helpers
* [ ] Request validation
* [ ] User registration
* [ ] User login
* [ ] JWT authentication
* [ ] Refresh token rotation
* [ ] Logout functionality
* [ ] Habit CRUD
* [ ] Daily check-ins
* [ ] Streak calculation
* [ ] Middleware
* [ ] Docker support
* [ ] Database migrations
* [ ] Unit tests

---

# Why This Project?

This project is being built to learn:

* Go fundamentals
* net/http
* PostgreSQL
* Contexts
* Middleware
* Authentication
* JWT
* Dependency Injection
* Repository Pattern
* Service Layer
* Graceful Shutdown
* Production-grade API design

---

# License

MIT
