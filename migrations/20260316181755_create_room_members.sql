-- +goose Up
CREATE TABLE IF NOT EXISTS room_members(
    room_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_message_read_id BIGINT DEFAULT 0,
    
    CONSTRAINT pk_room_members PRIMARY KEY (room_id, user_id),
    CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_room_members_user_id ON room_members(user_id);
CREATE INDEX IF NOT EXISTS idx_room_members_read ON room_members(room_id, last_message_read_id);

-- +goose Down
DROP TABLE IF EXISTS room_members;
