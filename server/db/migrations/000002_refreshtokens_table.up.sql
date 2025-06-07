CREATE TABLE refresh_tokens (
    user_id INT REFERENCES users(id),
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);