# 🔐 Advanced Job Scheduler & Auth Service

A robust, production-ready microservice architecture engineered in Go. This system combines secure authentication with an asynchronous job scheduling engine using Redis and background workers.

## ✨ Key Features & System Design

### 🛡️ State-of-the-Art Security
*   **Asymmetric Session Management:** Utilizes stateless, short-lived JWTs (Access Tokens) for rapid authorization and stateful, opaque Refresh Tokens for secure re-authentication.
*   **Cryptographic Hashing at Rest:** Refresh tokens are hashed using `SHA-256` before database persistence. Passwords are salted and hashed using `bcrypt`.
*   **Multi-Device Session Control:** Supports concurrent logins with a strict **5-session limit per user**, implementing an automated FIFO eviction strategy.

### ⚙️ Asynchronous Job Processing
*   **Redis-Backed Queue:** Implements a reliable producer-consumer pattern using Redis lists for high-throughput job handling.
*   **Decoupled Workers:** Background workers run in separate containers, ensuring that heavy tasks (like email delivery) don't block the main API response.
*   **SMTP Integration:** Robust email delivery system with proper header management and error handling.

### 🏗️ Architecture & Engineering Quality
*   **Layered Architecture:** Strict separation of concerns using a `Handler → Service → Store` pattern.
*   **Dockerized Infrastructure:** Fully containerized setup for Postgres, Redis, the API, and the Worker.
*   **Zero-Downtime Migrations:** Database schema versioning managed via `golang-migrate`.

### 🚀 Tech Stack
*   **Language:** Go 1.25
*   **Routing:** [chi](https://github.com/go-chi/chi)
*   **Database:** PostgreSQL 15 & Redis 7
*   **Messaging:** Redis Pub/Sub & Lists
*   **Infrastructure:** Docker & Docker Compose

---

## 📂 Project Structure

```text
cmd/
├── api/                    → HTTP Transport & Controllers
│   ├── main.go             → API Entrypoint
│   ├── app.go              → Server bootstrap
│   └── jobHandler.go       → Job submission logic
└── worker/                 → Background Worker
    ├── main.go             → Worker Entrypoint
    └── mail_sender.go      → Email processing logic

internals/                  → Core Domain Logic
├── db/                     → Postgres connection management
├── redis/                  → Redis client & queue utilities
├── jwt/                    → Token generation & validation
└── store/                  → Data Access Objects (DAO)

migrations/                 → SQL migration files (Up/Down)
```

---

## 🚀 Quick Start

### 1. Configure Environment
Create a `.env` file in the root directory:
```env
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password
JWT_SECRET=your-secure-secret
```

### 2. Run via Docker
```bash
docker compose up --build
```

### 3. Run Migrations
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```

---

## 📡 API Reference

### 1. Job Submission
`POST /jobs`
Enqueues a background job (e.g., sending a welcome email).

```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "name": "John Doe"}'
```

### 2. Authentication (Auth)
*   `POST /signup` - Create a new account.
*   `POST /login` - Authenticate and get tokens.
*   `POST /refreshToken` - Rotate your session.
*   `GET /health` - Protected health check.

---

## 🗄️ Database Schema

### Users Table
Stores hashed credentials and timestamps.

### Refresh Tokens Table
Stores hashed refresh tokens with a reference to the user, supporting the 5-session limit via FIFO eviction.
