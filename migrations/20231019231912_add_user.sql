-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    "uuid"       uuid primary key,
    phoneNumber  varchar(20) not null,
    firstName    varchar(64) not null,
    lastName     varchar(64),
    registeredAt timestamp without time zone
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table users;
-- +goose StatementEnd
