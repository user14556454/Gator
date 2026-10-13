-- +goose Up   

CREATE TABLE users(
    id uuid PRIMARY KEY,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;    