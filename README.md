# 🔐 Advanced Authentication Service in Go

A robust, production-ready authentication microservice engineered in Go. This service implements modern security best practices, scalable session management, and a clean, layered architecture designed for high maintainability.

## ✨ Key Features & System Design

### 🛡️ State-of-the-Art Security
*   **Asymmetric Session Management:** Utilizes stateless, short-lived JWTs (Access Tokens) for rapid authorization and stateful, opaque Refresh Tokens for secure re-authentication.
*   **Cryptographic Hashing at Rest:** Refresh tokens are hashed using `SHA-256` before database persistence, preventing token hijacking in the event of a database compromise. Passwords are salted and hashed using `bcrypt`.
*   **Multi-Device Session Control:** Supports concurrent logins with a strict **5-session limit per user**. Implements an automated First-In-First-Out (FIFO) eviction strategy via SQL offset queries to gracefully log out the oldest device when the limit is breached.

### 🏗️ Architecture & Engineering Quality
*   **Layered Architecture:** Strict separation of concerns using a `Handler → Service → Store` pattern, allowing for modularity and isolation of business logic from transport and storage layers.
*   **Interface-Driven Development:** Employs Go's implicit interfaces (Duck Typing) for data repositories (`UserRepository`, `TokenService`). This enables robust dependency injection and trivial mocking for unit tests.
*   **Context Propagation:** End-to-end `context.Context` passing from HTTP handlers down to the database layer, ensuring graceful degradation, timeout management, and request cancellation.
*   **Zero-Downtime Migrations:** Database schema versioning managed via `golang-migrate`, ensuring deterministic, repeatable schema evolutions.

### 🚀 Tech Stack
*   **Language:** Go 1.25
*   **Routing:** [chi](https://github.com/go-chi/chi) (Lightweight, idiomatic router)
*   **Database:** PostgreSQL 15
*   **Cryptography:** `golang-jwt/jwt/v5`, `x/crypto/bcrypt`, `crypto/sha256`
*   **Infrastructure:** Docker & Docker Compose

---

## 📂 Project Structure

```text
cmd/                        → Application entrypoint & HTTP transport
├── main.go                 → Application bootstrap (Config, DB conn, Dependency Injection)
├── app.go                  → Router configuration and server lifecycle
├── auth_handler.go         → Login, Signup, and Middleware implementations
├── refresh_token_handler.go→ Secure token rotation logic
└── health.go               → Readiness/Liveness probes

internals/                  → Core domain logic and infrastructure
├── db/
│   └── db.go               → PostgreSQL connection pool management
├── jwt/
│   └── jwt.go              → Cryptographic token generation, validation, and hashing
└── store/
    ├── store.go            → Data access object (DAO) aggregator
    ├── user.go             → User repository (CRUD operations)
    └── refresh_token.go    → Session management repository (Storage & Eviction)

migrations/                 → SQL migration files (Up/Down)
```

---

## 🚀 Quick Start

### Prerequisites
- [Go 1.25+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Run via Docker (Recommended)
Bootstraps the entire environment, including the PostgreSQL container and the Go API.

```bash
docker compose up --build
```
*In a separate terminal, run migrations:*
```bash
migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```

---

## 📡 API Reference

### 1. Create Account
`POST /signup`
Creates a new user and provisions their initial session.

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username": "engineer", "password": "securepassword123"}'
```
**Response (201 Created):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "uYxg8v3..."
}
```

### 2. Authenticate
`POST /login`
Verifies credentials and provisions a new device session (respecting the 5-device limit).

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username": "engineer", "password": "securepassword123"}'
```
**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "pLx4s1..."
}
```

### 3. Rotate Session
`POST /refreshToken`
Exchanges a valid, plaintext refresh token for a new short-lived access token. The server hashes the provided token and verifies it against the database.

```bash
curl -X POST http://localhost:8080/refreshToken \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "pLx4s1..."}'
```
**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "pLx4s1..."
}
```

### 4. Protected Route (Health)
`GET /health`
Validates the JWT signature and expiration.

```bash
curl -X GET http://localhost:8080/health \
  -H "Authorization: Bearer <access_token>"
```
**Response (200 OK):** `Good health`

---

## 🗄️ Database Schema

```sql
CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id TEXT REFERENCES users(username),
    token TEXT UNIQUE,       -- Stores the SHA-256 hash of the token
    expires_at TIMESTAMP
);
```

