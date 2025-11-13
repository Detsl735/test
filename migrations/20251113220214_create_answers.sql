-- +goose Up
-- +goose StatementBegin
CREATE TABLE answers (
    id           SERIAL PRIMARY KEY,
    question_id  INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    user_id      VARCHAR(64) NOT NULL,
    text         TEXT NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
