-- +goose Up
-- +goose StatementBegin
CREATE TABLE event_favorites
(
    userUuid  uuid      not null references users (uuid),
    eventUuid uuid      not null references events (uuid),
    addedAt   timestamp not null default current_timestamp,
    primary key (userUuid, eventUuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE event_favorites;
-- +goose StatementEnd
