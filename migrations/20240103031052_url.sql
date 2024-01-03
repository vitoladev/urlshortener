-- +goose Up
-- +goose StatementBegin
CREATE TABLE url (
    id SERIAL PRIMARY KEY,
    original_url VARCHAR(255) NOT NULL,
    short_url VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS url;
-- +goose StatementEnd
