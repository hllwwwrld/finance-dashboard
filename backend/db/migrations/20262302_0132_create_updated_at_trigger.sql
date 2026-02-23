-- +goose Up
-- +goose StatementBegin
CREATE
    OR REPLACE FUNCTION set_updated_at_column() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now() AT TIME ZONE 'utc';
    RETURN NEW;
END;
$$
    language 'plpgsql';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd