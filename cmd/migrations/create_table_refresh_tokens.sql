CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id UUID PRIMARY KEY,
    token TEXT NOT NULL
);