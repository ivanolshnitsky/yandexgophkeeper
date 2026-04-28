CREATE TABLE IF NOT EXISTS data (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,
    value JSONB NOT NULL,
    meta TEXT
);