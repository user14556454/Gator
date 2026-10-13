-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at timestamp without time zone;

-- +goose Down
ALTER TABLE feeds
DROP COLUMN last_fetched_at;