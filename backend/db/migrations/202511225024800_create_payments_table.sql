-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.payment
(
    id text DEFAULT gen_random_uuid() PRIMARY KEY,
    "name" text DEFAULT '' NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment
-- +goose StatementEnd