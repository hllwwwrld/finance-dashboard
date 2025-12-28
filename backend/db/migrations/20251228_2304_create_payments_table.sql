-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    id text DEFAULT gen_random_uuid() PRIMARY KEY,
    monthly_income int DEFAULT 0,
    created_at TIMESTAMP WITH TIMEZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIMEZONE DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd