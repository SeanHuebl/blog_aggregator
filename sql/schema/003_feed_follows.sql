-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID,
    feed_id UUID,
    CONSTRAINT user_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT feed_fk FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);
-- +goose Down
DROP TABLE feed_follows CASCADE;