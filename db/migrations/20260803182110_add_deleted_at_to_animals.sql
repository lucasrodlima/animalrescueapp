-- +goose Up
ALTER TABLE animals
ADD COLUMN deleted_at TIMESTAMP NULL;

-- +goose Down
ALTER TABLE animals
DROP COLUMN deleted_at;
