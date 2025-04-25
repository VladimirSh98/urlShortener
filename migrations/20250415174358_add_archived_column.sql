-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls ADD COLUMN archived bool DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls DROP COLUMN IF EXISTS archived;
-- +goose StatementEnd
