-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER user_updated_at
    BEFORE UPDATE
    ON "user"
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SEQUENCE user_updated_at;
-- +goose StatementEnd