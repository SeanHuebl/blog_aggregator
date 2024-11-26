-- +goose Up
-- Add the `last_fetched_at` column to track the last time an RSS feed was fetched
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;
-- +goose Down
-- Remove the `last_fetched_at` column from the `feeds` table
ALTER TABLE feeds DROP COLUMN last_fetched_at CASCADE;