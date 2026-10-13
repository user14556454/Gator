-- +goose Up
CREATE TABLE feed_follows (
    id uuid PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    user_id uuid NOT NULL,
    feed_id uuid NOT NULL,
    CONSTRAINT fk_users FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_feeds FOREIGN KEY(feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT uq_feed_user UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;