CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id TEXT REFERENCES users(username),
    token TEXT UNIQUE,
    expires_at TIMESTAMP
);