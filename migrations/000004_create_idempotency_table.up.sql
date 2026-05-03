CREATE TABLE idempotency (
    id SERIAL PRIMARY KEY,
    idempotency_key varchar(255) NOT NULL UNIQUE,
    created_at timestamp default now()
);