CREATE TABLE IF NOT EXISTS data (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value TEXT NOT NULL,
    meta TEXT,
    updated_at TIMESTAMP NOT NULL
);