-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment
(
    id TEXT DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id TEXT NOT NULL,
    "name" TEXT NOT NULL,
    amount INT NOT NULL,
    due_day INT NOT NULL,
    category TEXT NOT NULL,
    color TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment
-- +goose StatementEnd