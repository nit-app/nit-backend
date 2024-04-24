-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN isAdmin bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN isAdmin;
-- +goose StatementEnd
