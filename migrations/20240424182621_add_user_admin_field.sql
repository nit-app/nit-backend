-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN isAdmin bool not null default false;

INSERT INTO users (uuid, phonenumber, firstname, lastname, registeredat, isAdmin)
values ('72abf5a8-0400-11ef-829d-dca6329addd6', '78889992255', 'Egor', 'Krid', '2020-01-01T20:05:02', true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM users
WHERE uuid = '72abf5a8-0400-11ef-829d-dca6329addd6';

ALTER TABLE users
    DROP COLUMN isAdmin;
-- +goose StatementEnd
