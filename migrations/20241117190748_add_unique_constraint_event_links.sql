-- +goose Up
-- +goose StatementBegin
ALTER TABLE event_external_links
    ADD CONSTRAINT unique_event_url UNIQUE ("eventUuid", url);

ALTER TABLE event_external_links
    ADD CONSTRAINT unique_event_title UNIQUE ("eventUuid", title);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event_external_links
DROP CONSTRAINT unique_event_url;

ALTER TABLE event_external_links
DROP CONSTRAINT unique_event_title;
-- +goose StatementEnd
