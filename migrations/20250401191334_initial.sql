-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id varchar NOT NULL,
    original_url varchar NOT NULL,
    created_at timestamp default NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
