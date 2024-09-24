CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    token  TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    UNIQUE (user_id),
);