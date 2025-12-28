-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.payments
(
    id TEXT DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id TEXT NOT NULL,
    "name" TEXT NOT NULL,
    amount INT NOT NULL,
    due_date date NOT NULL,
    category TEXT NOT NULL,
    color TEXT NOT NULL,
    created_at TIMESTAMP WITH TIMEZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIMEZONE DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments
-- +goose StatementEnd