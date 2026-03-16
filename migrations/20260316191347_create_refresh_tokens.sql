-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_tokens(
    user_id BIGINT NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT pk_tokens PRIMARY KEY (user_id, token)
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
