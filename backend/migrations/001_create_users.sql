CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    nme VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at UNIXTIMESTAMP NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at UNIXTIMESTAMP NOT NULL DEFAULT (strftime('%s', 'now'))
)