-- +goose Up
-- +goose StatementBegin
SELECT 'CREATE EXTENSION IF NOT EXISTS postgis;';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'DROP EXTENSION IF EXISTS postgis';
-- +goose StatementEnd
