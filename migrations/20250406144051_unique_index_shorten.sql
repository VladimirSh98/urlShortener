-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS urls_original_url_udx ON urls(original_url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX urls_original_url_udx;
-- +goose StatementEnd
