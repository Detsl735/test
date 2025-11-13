-- +goose Up
-- +goose StatementBegin
CREATE TABLE questions (
    id          SERIAL PRIMARY KEY,
    text        TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS questions;
-- +goose StatementEnd
