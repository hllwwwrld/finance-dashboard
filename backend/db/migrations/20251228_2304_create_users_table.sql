-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    id text DEFAULT gen_random_uuid() PRIMARY KEY,
    login TEXT NOT NULL,
    password TEXT NOT NULL,
    monthly_income INT DEFAULT 0,
    created_at TIMESTAMP WITH TIMEZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIMEZONE DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd